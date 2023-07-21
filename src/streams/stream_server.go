package streams

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/services"
)

var PingPacket = []byte{0x50, 0x49, 0x4e, 0x47}
var PongPacket = []byte{0x50, 0x4f, 0x4e, 0x47}
var DataPacket = []byte{0x44, 0x41, 0x54, 0x41}

type StreamState struct {
	room *dto.Room
	mu   sync.RWMutex

	authorState       *User
	guestState        *User
	authorStateServer Stream_StreamStateServer
	guestStateServer  Stream_StreamStateServer

	authorStream       *websocket.Conn
	guestStream        *websocket.Conn
	authorStreamServer Stream_AVStreamServer
	guestStreamServer  Stream_AVStreamServer
}

// StreamServiceServer : we will use websockets to accept stream from user
// because stream server on browser is not working
type StreamServiceServer struct {
	UnimplementedStreamServer
	roomService *services.RoomService
	authService *services.AuthService

	addr string

	wg           *sync.WaitGroup
	mu           sync.RWMutex
	streamStates map[uuid.UUID]*StreamState
}

func (s *StreamServiceServer) Init(addr string, wg *sync.WaitGroup) error {
	s.addr = addr
	s.wg = wg

	return nil
}

func (s *StreamServiceServer) Run(ctx context.Context) error {
	s.wg.Add(1)

	log.Println("StreamServiceServer:: starting")

	// we create grpc without tls, envoy will terminate it
	grpcServer := grpc.NewServer()
	RegisterStreamServer(grpcServer, s)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errors.New(fmt.Sprintf("StreamServiceServer:: [grpc] failed to listen: %v", err))
	}

	// serve normal grpc
	go func() {
		err = grpcServer.Serve(ln)
		if err != nil {
			log.Fatalf("StreamServiceServer:: [grpc] unhandled error: %v", err)
		}
	}()
	log.Println("StreamServiceServer:: [grpc] started")

	// serve basic ws
	mux := http.NewServeMux()
	mux.Handle("/Streams/Connect/", websocket.Handler(s.handleWS))
	server := http.Server{
		Addr:    s.addr + "0",
		Handler: mux,
	}
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("StreamServiceServer:: [ws] unhandled error: %v", err)
		}
	}()

	log.Println("StreamServiceServer:: [ws] started")

	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = server.Shutdown(shutdownContext)
		if err != nil {
			log.Println("StreamServiceServer:: graceful shutdown error: ", err)
		}
		grpcServer.Stop()
		log.Println("StreamServiceServer:: shutdown complete")
	}()
	return nil
}

func (s *StreamServiceServer) StreamState(userRequest *User, stream Stream_StreamStateServer) error {
	user, err := s.authService.GetUserByString(userRequest.UserUUID)
	if err != nil {
		return err
	}
	room := s.roomService.State(user)
	if room != nil {
		state := s.getState(room)
		state.mu.Lock()
		if room.Author == user {
			state.authorState = userRequest
			state.authorStateServer = stream
		} else {
			state.guestState = userRequest
			state.authorStateServer = stream
		}
		state.mu.Unlock()
		s.sendStateUpdates(state)

		<-stream.Context().Done()
	}
	return nil
}

func (s *StreamServiceServer) ChangeState(ctx context.Context, userRequest *User) (*Ack, error) {
	user, err := s.authService.GetUserByString(userRequest.UserUUID)
	if err != nil {
		return nil, err
	}
	room := s.roomService.State(user)
	if room != nil {
		state := s.getState(room)
		state.mu.Lock()
		if room.Author == user {
			state.authorState = userRequest
		} else {
			state.guestState = userRequest
		}
		state.mu.Unlock()
		s.sendStateUpdates(state)
	}
	return &Ack{}, nil
}

func (s *StreamServiceServer) AVStream(userRequest *User, stream Stream_AVStreamServer) error {
	user, err := s.authService.GetUserByString(userRequest.UserUUID)
	if err != nil {
		return err
	}
	room := s.roomService.State(user)
	if room != nil {
		state := s.getState(room)
		state.mu.Lock()
		if room.Author == user {
			state.authorState = userRequest
			state.authorStreamServer = stream
		} else {
			state.guestState = userRequest
			state.guestStreamServer = stream
		}
		state.mu.Unlock()

		s.sendStateUpdates(state)

		<-stream.Context().Done()
	}
	return nil
}

func (s *StreamServiceServer) handleWS(ws *websocket.Conn) {
	userUUID := string([]byte(ws.Request().URL.Path)[17:53])
	roomUUID := string([]byte(ws.Request().URL.Path)[54:])
	user, room, err := s.userAndRoom(userUUID, roomUUID)
	if err != nil {
		ws.Close()
	} else {
		if room.Guest == user || room.Author == user {
			go s.readLoop(ws, user, room)
		} else {
			ws.Close()
		}
	}
}

func (s *StreamServiceServer) readLoop(ws *websocket.Conn, user *dto.User, room *dto.Room) {
	// buffer for reads, 128kb per connection
	buf := make([]byte, 1024*128)
	var isAuthor bool

	// stream state fetch / creation

	state := s.getState(room)
	state.mu.Lock()
	if room.Author == user {
		state.authorStream = ws
		isAuthor = true
	} else {
		state.guestStream = ws
		isAuthor = false
	}
	state.mu.Unlock()

	for {
		length, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("StreamServiceServer:: [ws] read error")
			continue
		}
		// we handle ws connections alive on front-end
		if bytes.Equal(PingPacket, buf[:3]) {
			_, err := ws.Write(PongPacket)
			if err != nil {
				log.Println("StreamServiceServer:: [ws] write error")
				ws.Close()
				break
			}
			continue
		}
		// AV Data came need resend to grpc web stream
		if bytes.Equal(DataPacket, buf[:3]) {
			var stream Stream_AVStreamServer
			if isAuthor && state.guestStreamServer != nil {
				stream = state.guestStreamServer
			}
			if !isAuthor && state.authorStreamServer != nil {
				stream = state.authorStreamServer
			}
			if stream != nil {
				// we ignore errors, it's safe
				_ = state.guestStreamServer.Send(&AVFrameData{
					UserUUID:  user.UUID.String(),
					FrameData: buf[4:length],
				})
			}
			continue
		}
		log.Println("StreamServiceServer:: [ws] unhandled message: ", string(buf[:length]))
	}
}

func (s *StreamServiceServer) createState(room *dto.Room) *StreamState {
	roomState := &StreamState{
		room:               room,
		mu:                 sync.RWMutex{},
		authorStream:       nil,
		guestStream:        nil,
		authorStreamServer: nil,
		guestStreamServer:  nil,
	}
	s.streamStates[room.UUID] = roomState
	return roomState
}

func (s *StreamServiceServer) getState(room *dto.Room) *StreamState {
	s.mu.Lock()
	state, exists := s.streamStates[room.UUID]
	if !exists {
		state = s.createState(room)
	}
	s.mu.Unlock()
	return state
}

func (s *StreamServiceServer) sendStateUpdates(state *StreamState) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	stateMessage := &StateMessage{
		Time:   time.Now().Unix(),
		UUID:   uuid.NewString(),
		Author: state.authorState,
		Guest:  state.guestState,
	}
	if state.authorStateServer != nil {
		state.authorStateServer.Send(stateMessage)
	}
	if state.guestStateServer != nil {
		state.guestStateServer.Send(stateMessage)
	}
}

func (s *StreamServiceServer) userAndRoom(userStringUUID string, roomStringUUID string) (*dto.User, *dto.Room, error) {
	user, err := s.authService.GetUserByString(userStringUUID)
	if err != nil {
		return nil, nil, err
	}
	room, err := s.roomService.GetRoom(user, roomStringUUID)
	if err != nil {
		return nil, nil, err
	}
	return user, room, nil
}

func NewStreamServiceServer(roomService *services.RoomService, authService *services.AuthService) *StreamServiceServer {
	return &StreamServiceServer{
		roomService:  roomService,
		authService:  authService,
		streamStates: make(map[uuid.UUID]*StreamState),
	}
}
