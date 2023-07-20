package streams

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/services"
)

type StreamState struct {
	room               *dto.Room
	mu                 sync.RWMutex
	authorStreamClient Stream_AVStreamClient
	authorStreamServer Stream_AVStreamServer
	guestStreamClient  Stream_AVStreamClient
	guestStreamServer  Stream_AVStreamServer
}

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

func (s *StreamServiceServer) Run(ctx context.Context) {
	s.wg.Add(1)

	log.Println("StreamServiceServer:: starting")

	// we create grpc without tls, envoy will terminate it
	grpcServer := grpc.NewServer()
	RegisterStreamServer(grpcServer, s)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// serve normal grpc
	go grpcServer.Serve(ln)
	log.Println("StreamServiceServer:: grpc started")

	go func() {
		defer s.wg.Done()

		<-ctx.Done()

		grpcServer.Stop()

		log.Println("StreamServiceServer:: shutdown complete")
	}()
}

func NewStreamServiceServer(roomService *services.RoomService, authService *services.AuthService) *StreamServiceServer {
	return &StreamServiceServer{
		roomService:  roomService,
		authService:  authService,
		streamStates: make(map[uuid.UUID]*StreamState),
	}
}
