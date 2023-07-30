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
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
)

type ChatServer struct {
	stateProvider *RoomStateProvider

	addr  string
	wg    *sync.WaitGroup
	mu    sync.RWMutex
	chats map[uuid.UUID]*ChatState
	chat.UnimplementedChatServer
}

type ChatState struct {
	room     *dto.Room
	messages []*chat.ChatMessage
	msgChan  chan *chat.ChatMessage
	mu       sync.RWMutex
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
	state, _, err := s.stateProvider.GetByUserAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		fmt.Println("Error on fetching user and room, room might be already closed or never exists.")
		return nil, err
	}
	state.chat.mu.RLock()
	defer state.chat.mu.RUnlock()
	return &chat.HistoryResponse{
		Messages: state.chat.messages,
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, request *chat.SendMessageRequest) (*chat.Empty, error) {
	// we must check User and Room correctness
	state, user, err := s.stateProvider.GetByUserAndRoom(
		request.GetAuthData().GetUUID(),
		request.GetAuthData().GetChatUUID(),
	)
	if err != nil {
		return nil, err
	}
	if request.GetMsg() == "" {
		return &chat.Empty{}, nil
	}
	// all fine, lets proceed
	chatMessage := &chat.ChatMessage{
		UUID:     uuid.New().String(),
		UserUUID: user.UUID.String(),
		UserName: user.Name,
		Time:     time.Now().Unix(),
		Msg:      request.Msg,
	}

	state.chat.msgChan <- chatMessage

	return &chat.Empty{}, nil
}

func (s *ChatServer) Listen(request *chat.AuthRequest, stream chat.Chat_ListenServer) error {
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
		stream.Send(&chat.ChatMessage{
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

func NewChatServer(provider *RoomStateProvider) *ChatServer {
	return &ChatServer{
		stateProvider: provider,
	}
}

func AddChatState(roomState *RoomState) {
	chatState := &ChatState{
		room: roomState.room,
		messages: []*chat.ChatMessage{{
			UUID:     uuid.NewString(),
			UserUUID: uuid.NewString(),
			UserName: "Server",
			Time:     0,
			Msg:      "Welcome to chat",
		}},
		msgChan: make(chan *chat.ChatMessage, 4),
	}
	go func(state *ChatState) {
		for {
			message, ok := <-state.msgChan
			if !ok {
				return
			}
			state.mu.Lock()
			state.messages = append(state.messages, message)
			state.mu.Unlock()
			if roomState.author.chatStream.stream != nil {
				roomState.author.chatStream.stream.Send(message)
			}
			if roomState.guest != nil && roomState.guest.chatStream.stream != nil {
				roomState.guest.chatStream.stream.Send(message)
			}
		}
	}(chatState)
	roomState.chat = chatState
}
