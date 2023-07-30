package inroom

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
	"web-grpc-video-chat/src/inroom/stream"
)

var PingPacket = []byte{0x50, 0x49, 0x4e, 0x47}
var PongPacket = []byte{0x50, 0x4f, 0x4e, 0x47}
var DataPacket = []byte{0x44, 0x41, 0x54, 0x41}

// StreamServer : we will use websockets to accept stream from user
// because stream server on browser is not working
type StreamServer struct {
	wg            *sync.WaitGroup
	stateProvider *RoomStateProvider
	addr          string
	stream.UnimplementedStreamServer
}

func (s *StreamServer) Init(addr string, wg *sync.WaitGroup) error {
	s.addr = addr
	s.wg = wg

	return nil
}

func (s *StreamServer) Run(ctx context.Context) error {
	s.wg.Add(1)

	log.Println("StreamServiceServer:: starting")

	// we create grpc without tls, envoy will terminate it
	grpcServer := grpc.NewServer()
	stream.RegisterStreamServer(grpcServer, s)

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

	// Websockets listen and serve
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

func (s *StreamServer) ChangeState(ctx context.Context, userRequest *stream.User) (*stream.Ack, error) {
	state, user, err := s.stateProvider.GetByUserAndRoom(userRequest.GetUserUUID(), userRequest.GetUserRoom())
	if err != nil {
		return nil, err
	}

	state.UpdateUserState(user, userRequest)
	s.SendStateUpdate(state)

	return &stream.Ack{}, nil
}

func (s *StreamServer) StreamState(userRequest *stream.User, stream stream.Stream_StreamStateServer) error {
	state, user, err := s.stateProvider.GetByUserAndRoom(userRequest.GetUserUUID(), userRequest.GetUserRoom())
	if err != nil {
		return err
	}

	closeConnCh, err := state.StateStreamConnect(user, stream)
	if err != nil {
		return err
	}

	state.UpdateUserState(user, userRequest)
	s.SendStateUpdate(state)

	select {
	case <-closeConnCh:
	case <-stream.Context().Done():
	case <-state.roomCtx.Done():
	}

	return nil
}

func (s *StreamServer) AVStream(userRequest *stream.User, stream stream.Stream_AVStreamServer) error {
	state, user, err := s.stateProvider.GetByUserAndRoom(userRequest.GetUserUUID(), userRequest.GetUserRoom())
	if err != nil {
		return err
	}

	closeConnCh, err := state.AVStreamConnect(user, stream)
	if err != nil {
		return err
	}

	state.UpdateUserState(user, userRequest)
	s.SendStateUpdate(state)

	select {
	case <-closeConnCh:
	case <-stream.Context().Done():
	case <-state.roomCtx.Done():
	}
	return nil
}

func (s *StreamServer) SendStateUpdate(state *RoomState) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	stateMessage := &stream.StateMessage{
		Time:   time.Now().Unix(),
		UUID:   uuid.NewString(),
		Author: state.author.state,
		Guest:  nil,
	}
	if state.guest != nil {
		stateMessage.Guest = state.guest.state
	}
	if state.author.stateStream.stream != nil {
		state.author.stateStream.stream.Send(stateMessage)
	}
	if state.guest != nil && state.guest.stateStream.stream != nil {
		state.guest.stateStream.stream.Send(stateMessage)
	}
}

func (s *StreamServer) handleWS(ws *websocket.Conn) {
	userUUID := string([]byte(ws.Request().URL.Path)[17:53])
	roomUUID := string([]byte(ws.Request().URL.Path)[54:])
	state, user, err := s.stateProvider.GetByUserAndRoom(userUUID, roomUUID)
	if err != nil {
		ws.Close()
	} else {
		if state.author.user == user || state.guest != nil && state.guest.user == user {
			go s.readLoop(ws, state, user)
		} else {
			ws.Close()
		}
	}
}

func (s *StreamServer) readLoop(ws *websocket.Conn, state *RoomState, user *dto.User) {
	// buffer for reads, 128kb per connection
	buf := make([]byte, 1024*128)

	// can we fetch userState
	userState := state.GetUserState(user)
	if userState == nil {
		ws.Close()
		return
	}

	// some bad design magic, can be really improved much better
	state.mu.Lock()
	if userState.inputStream != nil {
		userState.inputStream.Close()
	}
	userState.inputStream = ws
	state.mu.Unlock()

	go func() {
		<-state.roomCtx.Done()
		ws.Close()
	}()

	for {
		length, err := ws.Read(buf)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println("StreamServiceServer:: [ws] read error", err)
			continue
		}
		// we handle ws connections alive on front-end
		if bytes.Equal(PingPacket, buf[:3]) {
			respondPingPing(ws)
			continue
		}
		// AV Data came need resend to grpc web stream
		if bytes.Equal(DataPacket, buf[:3]) {
			handleDataPacket(buf, length, state, userState)
			continue
		}
		log.Println("StreamServiceServer:: [ws] unhandled message: ", string(buf[:length]))
	}
}

func handleDataPacket(buf []byte, length int, state *RoomState, userState *UserState) {
	avStream := state.GetOpponentDataStream(userState)
	// we ignore errors, it's safe
	_ = avStream.Send(&stream.AVFrameData{
		UserUUID:  userState.user.UUID.String(),
		FrameData: buf[4:length],
	})
}

func respondPingPing(ws *websocket.Conn) {
	_, err := ws.Write(PongPacket)
	if err != nil {
		log.Println("StreamServiceServer:: [ws] write error")
		ws.Close()
	}
}

func NewStreamServer(
	provider *RoomStateProvider,
) *StreamServer {
	return &StreamServer{
		stateProvider: provider,
	}
}
