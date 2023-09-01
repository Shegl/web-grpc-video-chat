package pb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/internal/core/repo"
	"web-grpc-video-chat/src/internal/core/services"
	"web-grpc-video-chat/src/pb/chat"
)

type ChatListenServerAdapter struct {
	server chat.Chat_ListenServer
}

func (s *ChatListenServerAdapter) Send(msg domain.ChatMessage) error {
	return s.server.Send(buildPbChatMessage(msg))
}

type ChatServer struct {
	roomProvider *services.RoomManagerProvider
	repo         *repo.Repository
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
	// we must validate User and Room
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
	if manager.InRoom(user) {
		chatMessages := manager.GetChatHistory()
		pbMessages := make([]*chat.ChatMessage, len(chatMessages), len(chatMessages))
		for k, msg := range chatMessages {
			pbMessages[k] = buildPbChatMessage(msg)
		}
		return &chat.HistoryResponse{
			Messages: pbMessages,
		}, nil
	}
	return nil, errors.New("User not in room. ")
}

func (s *ChatServer) SendMessage(ctx context.Context, request *chat.SendMessageRequest) (v *chat.Empty, err error) {
	// we must validate User and Room
	if request.GetMsg() == "" {
		return
	}
	room, err := s.repo.FindRoomByString(request.GetAuthData().GetChatUUID())
	if err != nil {
		return
	}
	user, err := s.repo.FindUserByString(request.GetAuthData().GetUUID())
	if err != nil {
		return
	}
	manager := s.roomProvider.GetRoomManager(room)
	if manager != nil && manager.InRoom(user) {
		manager.ChatBroadcast(domain.ChatMessage{
			Time:     time.Now().Unix(),
			UUID:     uuid.New(),
			UserUUID: user.UUID,
			UserName: user.Name,
			Msg:      request.Msg,
		})
	}
	return
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
	if manager == nil || !manager.InRoom(user) {
		return errors.New("User not in room. ")
	}
	closeConnCh, err := manager.RoomChatConnect(user, &ChatListenServerAdapter{stream})
	if err != nil {
		return err
	}
	select {
	case <-closeConnCh:
	case <-stream.Context().Done():
	}
	return nil
}

func NewChatServer(provider *services.RoomManagerProvider, repo *repo.Repository) *ChatServer {
	return &ChatServer{
		roomProvider: provider,
		repo:         repo,
	}
}

func buildPbChatMessage(msg domain.ChatMessage) *chat.ChatMessage {
	return &chat.ChatMessage{
		UUID:     msg.GetUUID().String(),
		UserUUID: msg.GetUserUUID().String(),
		UserName: msg.GetUserName(),
		Time:     msg.GetTime(),
		Msg:      msg.GetMsg(),
	}
}
