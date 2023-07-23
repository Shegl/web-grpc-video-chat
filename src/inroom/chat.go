package inroom

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
	chat2 "web-grpc-video-chat/src/inroom/chat"
)

type ChatServer struct {
	stateProvider *RoomStateProvider

	addr  string
	wg    *sync.WaitGroup
	mu    sync.RWMutex
	chats map[uuid.UUID]*ChatState
	chat2.UnimplementedChatServer
}

func (s *ChatServer) Init(addr string, wg *sync.WaitGroup) error {
	s.addr = addr
	s.wg = wg

	return nil
}

func (s *ChatServer) Run(ctx context.Context) error {
	s.wg.Add(1)

	log.Println("ChatService:: starting")

	// we create grpc without tls, envoy will terminate it
	grpcServer := grpc.NewServer()
	chat2.RegisterChatServer(grpcServer, s)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("ChatServiceServer:: failed to listen: %v", err)
	}

	// serve normal grpc
	go grpcServer.Serve(ln)
	log.Println("ChatServiceServer:: grpc started")

	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		grpcServer.Stop()
		log.Println("ChatServiceServer:: shutdown complete")
	}()

	return nil
}

func (s *ChatServer) GetHistory(ctx context.Context, request *chat2.AuthRequest) (*chat2.HistoryResponse, error) {
	// we must check User and Room permissions
	state, _, err := s.stateProvider.GetByUserAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		fmt.Println("Error on fetching user and room, room might be already closed or never exists.")
		return nil, err
	}
	state.chat.mu.RLock()
	defer state.chat.mu.RUnlock()
	return &chat2.HistoryResponse{
		Messages: state.chat.messages,
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, request *chat2.SendMessageRequest) (*chat2.Empty, error) {
	// we must check User and Room correctness
	state, user, err := s.stateProvider.GetByUserAndRoom(
		request.GetAuthData().GetUUID(),
		request.GetAuthData().GetChatUUID(),
	)
	if err != nil {
		return nil, err
	}
	if request.GetMsg() == "" {
		return &chat2.Empty{}, nil
	}
	// all fine, lets proceed
	chatMessage := &chat2.ChatMessage{
		UUID:     uuid.New().String(),
		UserUUID: user.UUID.String(),
		UserName: user.Name,
		Time:     time.Now().Unix(),
		Msg:      request.Msg,
	}

	state.chat.msgChan <- chatMessage

	return &chat2.Empty{}, nil
}

func (s *ChatServer) Listen(request *chat2.AuthRequest, stream chat2.Chat_ListenServer) error {
	state, user, err := s.stateProvider.GetByUserAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		return err
	}
	closeConnCh, err := state.RoomChatConnect(user, stream)
	if err != nil {
		return err
	}

	select {
	case <-state.roomCtx.Done():
		stream.Send(&chat2.ChatMessage{
			UUID:     uuid.NewString(),
			UserUUID: uuid.NewString(),
			UserName: "Server",
			Time:     0,
			Msg:      "Room closed. You will be redirected. Bye! ",
		})
	case <-closeConnCh:
	case <-stream.Context().Done():
	}
	return nil
}

func NewChatServer(provider *RoomStateProvider) *ChatServer {
	return &ChatServer{
		stateProvider: provider,
	}
}
