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

type RoomStreamState struct {
	room *dto.Room
	mu   sync.RWMutex

	author *UserStreamState
	guest  *UserStreamState
}

type UserStreamState struct {
	state        *User
	stream       *websocket.Conn
	stateServer  Stream_StreamStateServer
	streamServer Stream_AVStreamServer
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
	streamStates map[uuid.UUID]*RoomStreamState
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

	// gRPC serve
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

	// Websockets listen and server
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("StreamServiceServer:: [ws] unhandled error: %v", err)
		}
	}()

	log.Println("StreamServiceServer:: [ws] started")

	// shutdown handler
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

	userState, roomState := s.getUserStreamState(user)
	if userState == nil {
		return nil
	}

	roomState.mu.Lock()
	userState.state = userRequest
	userState.stateServer = stream
	roomState.mu.Unlock()

	s.sendStateUpdates(roomState)

	<-stream.Context().Done()

	return nil
}

func (s *StreamServiceServer) ChangeState(ctx context.Context, userRequest *User) (*Ack, error) {
	user, err := s.authService.GetUserByString(userRequest.UserUUID)
	if err != nil {
		return nil, err
	}

	userState, roomState := s.getUserStreamState(user)
	if userState == nil {
		return &Ack{}, nil
	}

	roomState.mu.Lock()
	userState.state = userRequest
	roomState.mu.Unlock()

	s.sendStateUpdates(roomState)

	return &Ack{}, nil
}

func (s *StreamServiceServer) AVStream(userRequest *User, stream Stream_AVStreamServer) error {
	user, err := s.authService.GetUserByString(userRequest.UserUUID)
	if err != nil {
		return err
	}

	userState, roomState := s.getUserStreamState(user)
	if userState == nil {
		return nil
	}

	roomState.mu.Lock()
	userState.state = userRequest
	userState.streamServer = stream
	roomState.mu.Unlock()

	s.sendStateUpdates(roomState)

	<-stream.Context().Done()

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
	var opponent *UserStreamState

	// stream state fetch / creation

	userState, roomState := s.getUserStreamState(user)
	if userState == nil {
		ws.Close()
		return
	}

	roomState.mu.Lock()
	userState.stream = ws
	roomState.mu.Unlock()

	if userState == roomState.author {
		opponent = roomState.guest
	} else {
		opponent = roomState.author
	}

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
			if opponent.streamServer != nil {
				// we ignore errors, it's safe
				_ = opponent.streamServer.Send(&AVFrameData{
					UserUUID:  user.UUID.String(),
					FrameData: buf[4:length],
				})
			}
			continue
		}
		log.Println("StreamServiceServer:: [ws] unhandled message: ", string(buf[:length]))
	}
}

func (s *StreamServiceServer) createState(room *dto.Room) *RoomStreamState {
	roomState := &RoomStreamState{
		room: room,
		mu:   sync.RWMutex{},
		author: &UserStreamState{
			state:        nil,
			stream:       nil,
			stateServer:  nil,
			streamServer: nil,
		},
		guest: &UserStreamState{
			state:        nil,
			stream:       nil,
			stateServer:  nil,
			streamServer: nil,
		},
	}
	s.streamStates[room.UUID] = roomState
	return roomState
}

func (s *StreamServiceServer) getRoomState(room *dto.Room) *RoomStreamState {
	s.mu.Lock()
	roomState, exists := s.streamStates[room.UUID]
	if !exists {
		roomState = s.createState(room)
	}
	s.mu.Unlock()
	return roomState
}

func (s *StreamServiceServer) getUserStreamState(user *dto.User) (*UserStreamState, *RoomStreamState) {
	room := s.roomService.State(user)
	if room != nil {
		roomState := s.getRoomState(room)
		if roomState == nil {
			return nil, nil
		}
		roomState.mu.RLock()
		defer roomState.mu.RUnlock()
		var userState *UserStreamState
		if room.Author == user {
			userState = roomState.author
		} else {
			userState = roomState.guest
		}
		return userState, roomState
	}
	return nil, nil
}

func (s *StreamServiceServer) sendStateUpdates(roomState *RoomStreamState) {
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	stateMessage := &StateMessage{
		Time:   time.Now().Unix(),
		UUID:   uuid.NewString(),
		Author: roomState.author.state,
		Guest:  roomState.guest.state,
	}
	if roomState.author.stateServer != nil {
		roomState.author.stateServer.Send(stateMessage)
	}
	if roomState.guest.stateServer != nil {
		roomState.guest.stateServer.Send(stateMessage)
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
		streamStates: make(map[uuid.UUID]*RoomStreamState),
	}
}
