package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/services"
)

type ChatLog struct {
	room         *dto.Room
	messages     []*ChatMessage
	msgChan      chan *ChatMessage
	mu           sync.RWMutex
	authorStream Chat_ListenRequestServer
	guestStream  Chat_ListenRequestServer
}

type ChatServiceServer struct {
	UnimplementedChatServer
	roomService *services.RoomService
	authService *services.AuthService

	addr string

	wg       *sync.WaitGroup
	mu       sync.RWMutex
	chatLogs map[uuid.UUID]*ChatLog
}

func (s *ChatServiceServer) Init(addr string, wg *sync.WaitGroup) error {
	s.addr = addr
	s.wg = wg

	return nil
}

func (s *ChatServiceServer) Run(ctx context.Context) {
	s.wg.Add(1)

	log.Println("ChatServiceServer:: starting")

	// we create grpc without tsl, envoy will terminate it
	grpcServer := grpc.NewServer()
	RegisterChatServer(grpcServer, s)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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
}

func (s *ChatServiceServer) GetHistory(ctx context.Context, request *AuthRequest) (*HistoryResponse, error) {
	// we must check User and Room permissions
	_, room, err := s.userAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		fmt.Println("error on fetching user and room")
		return nil, err
	}
	return &HistoryResponse{
		Messages: s.getChat(room).messages,
	}, nil
}

func (s *ChatServiceServer) SendMessage(ctx context.Context, request *SendMessageRequest) (*Empty, error) {
	// we must check User and Room permissions
	user, room, err := s.userAndRoom(request.GetAuthData().GetUUID(), request.GetAuthData().GetChatUUID())
	if err != nil {
		return nil, err
	}
	if request.GetMsg() == "" {
		return &Empty{}, nil
	}
	// all fine, lets proceed
	chatMessage := &ChatMessage{
		UUID:     uuid.New().String(),
		UserUUID: user.UUID.String(),
		UserName: user.Name,
		Time:     time.Now().Unix(),
		Msg:      request.Msg,
	}
	chat := s.getChat(room)
	chat.msgChan <- chatMessage
	return &Empty{}, nil
}

func (s *ChatServiceServer) ListenRequest(request *AuthRequest, stream Chat_ListenRequestServer) error {
	user, room, err := s.userAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		return err
	}
	chat := s.getChat(room)
	chat.mu.Lock()
	if user == room.Author {
		chat.authorStream = stream
	} else {
		chat.guestStream = stream
	}
	chat.mu.Unlock()

	<-stream.Context().Done()

	return nil
}

func (s *ChatServiceServer) userAndRoom(userStringUUID string, chatStringUUID string) (*dto.User, *dto.Room, error) {
	user, err := s.getUser(userStringUUID)
	if err != nil {
		return nil, nil, err
	}
	// we must check is he in Room and have right to request history
	room, err := s.getRoom(user, chatStringUUID)
	if err != nil {
		return nil, nil, err
	}
	return user, room, nil
}

func (s *ChatServiceServer) getUser(stringUUID string) (*dto.User, error) {
	userUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	user, err := s.authService.GetUser(userUUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ChatServiceServer) getRoom(user *dto.User, stringUUID string) (*dto.Room, error) {
	roomUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	room := s.roomService.State(user)
	if room != nil && room.UUID == roomUUID {
		return room, nil
	}
	return nil, errors.New("Wrong room. ")
}

func (s *ChatServiceServer) getChat(room *dto.Room) *ChatLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	if chatLog, exists := s.chatLogs[room.UUID]; exists {
		return chatLog
	}
	chatLog := &ChatLog{
		room: room,
		messages: []*ChatMessage{{
			UUID:     "ajskdhfkjahsdf",
			UserUUID: "asldfjhlasjdf",
			UserName: "Bot",
			Time:     0,
			Msg:      "Welcome to chat",
		}},
		msgChan: make(chan *ChatMessage, 4),
	}
	go func() {
		for {
			message, ok := <-chatLog.msgChan
			if !ok {
				return
			}
			chatLog.messages = append(chatLog.messages, message)
			if chatLog.authorStream != nil {
				chatLog.authorStream.Send(message)
			}
			if chatLog.guestStream != nil {
				chatLog.authorStream.Send(message)
			}
		}
	}()
	s.chatLogs[room.UUID] = chatLog
	return chatLog
}

func NewChatServiceServer(roomService *services.RoomService, authService *services.AuthService) *ChatServiceServer {
	return &ChatServiceServer{
		roomService: roomService,
		authService: authService,
		chatLogs:    make(map[uuid.UUID]*ChatLog),
	}
}
