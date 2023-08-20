package inroom

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
)

type ChatServer struct {
	roomProvider *RoomProvider
	repo         *dto.Repository
	addr         string
	wg           *sync.WaitGroup
	mu           sync.RWMutex
	chat.UnimplementedChatServer
}

func (s *ChatServer) Init(addr string, wg *sync.WaitGroup) {
	s.addr = addr
	s.wg = wg
}

func (s *ChatServer) Run(ctx context.Context) error {
	s.wg.Add(1)
	log.Println("ChatService:: starting")

	// we create grpc without tls, envoy will terminate it
	grpcServer := grpc.NewServer()
	chat.RegisterChatServer(grpcServer, s)

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

func (s *ChatServer) GetHistory(ctx context.Context, request *chat.AuthRequest) (*chat.HistoryResponse, error) {
	// we must check User and Room permissions
	room, err := s.repo.FindRoomByString(request.GetChatUUID())
	if room == nil || err != nil {
		return nil, errors.New("No such room or uuid is not valid. ")
	}
	user, err := s.repo.FindUserByString(request.GetUUID())
	if err != nil {
		return nil, err
	}
	manager := s.roomProvider.GetRoomManager(room)
	if manager == nil {
		panic(errors.New("Something went wrong, not as designed. "))
	}
	if !manager.inRoom(user) {
		return nil, errors.New("User not in room. ")
	}
	return &chat.HistoryResponse{
		Messages: manager.getChatHistory(),
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, request *chat.SendMessageRequest) (*chat.Empty, error) {
	// we must check User and Room correctness
	if request.GetMsg() == "" {
		return &chat.Empty{}, nil
	}
	room, err := s.repo.FindRoomByString(request.GetAuthData().GetChatUUID())
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindUserByString(request.GetAuthData().GetUUID())
	if err != nil {
		return nil, err
	}
	manager := s.roomProvider.GetRoomManager(room)
	if manager != nil && manager.inRoom(user) {
		manager.chatBroadcast(&chat.ChatMessage{
			UUID:     uuid.New().String(),
			UserUUID: user.UUID.String(),
			UserName: user.Name,
			Time:     time.Now().Unix(),
			Msg:      request.Msg,
		})
	}
	return &chat.Empty{}, nil
}

func (s *ChatServer) Listen(request *chat.AuthRequest, stream chat.Chat_ListenServer) error {
	room, err := s.repo.FindRoomByString(request.GetChatUUID())
	if err != nil {
		return err
	}
	user, err := s.repo.FindUserByString(request.GetChatUUID())
	if err != nil {
		return err
	}
	manager := s.roomProvider.GetRoomManager(room)
	if manager == nil || !manager.inRoom(user) {
		return errors.New("User not in room. ")
	}
	closeConnCh, err := manager.roomChatConnect(user, stream)
	if err != nil {
		return err
	}
	select {
	case <-manager.roomCtx.Done():
		manager.chatBroadcast(&chat.ChatMessage{
			UUID:     uuid.NewString(),
			UserUUID: uuid.NewString(),
			UserName: "Server",
			Time:     0,
			Msg:      "Room closed by author or by server. You will be redirected. Bye! ",
		})
	case <-closeConnCh:
	case <-stream.Context().Done():
	}
	return nil
}

func NewChatServer(provider *RoomProvider, repo *dto.Repository) *ChatServer {
	return &ChatServer{
		roomProvider: provider,
		repo:         repo,
	}
}
