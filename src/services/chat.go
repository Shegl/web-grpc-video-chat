package services

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
	"web-grpc-video-chat/src/chat"
	"web-grpc-video-chat/src/dto"
)

type ChatState struct {
	room     *dto.Room
	messages []*chat.ChatMessage
	msgChan  chan *chat.ChatMessage
	mu       sync.RWMutex
}

type ChatService struct {
	chat.UnimplementedChatServer

	roomStateProvider *RoomStateProvider

	addr string

	wg    *sync.WaitGroup
	mu    sync.RWMutex
	chats map[uuid.UUID]*ChatState
}

func (s *ChatService) Init(addr string, wg *sync.WaitGroup) error {
	s.addr = addr
	s.wg = wg

	return nil
}

func (s *ChatService) Run(ctx context.Context) error {
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

func (s *ChatService) MakeChatState(room *dto.Room) *ChatState {
	chatState := &ChatState{
		room: room,
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
			roomState := s.roomStateProvider.GetRoomState(room)
			if roomState != nil {
				if roomState.author.chatStream.stream != nil {
					roomState.author.chatStream.stream.Send(message)
				}
				if roomState.guest != nil && roomState.guest.chatStream.stream != nil {
					roomState.guest.chatStream.stream.Send(message)
				}
			}
		}
	}(chatState)
	return chatState
}

func (s *ChatService) GetHistory(ctx context.Context, request *chat.AuthRequest) (*chat.HistoryResponse, error) {
	// we must check User and Room permissions
	state, _, err := s.stateByUserAndRoom(request.GetUUID(), request.GetChatUUID())
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

func (s *ChatService) SendMessage(ctx context.Context, request *chat.SendMessageRequest) (*chat.Empty, error) {
	// we must check User and Room correctness
	state, user, err := s.stateByUserAndRoom(request.GetAuthData().GetUUID(), request.GetAuthData().GetChatUUID())
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

func (s *ChatService) Listen(request *chat.AuthRequest, stream chat.Chat_ListenServer) error {
	state, user, err := s.stateByUserAndRoom(request.GetUUID(), request.GetChatUUID())
	if err != nil {
		return err
	}
	closeConnCh, err := state.RoomChatConnected(user, stream)
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
			Msg:      "Room closed. You will be redirected. Bye! ",
		})
	case <-closeConnCh:
	case <-stream.Context().Done():
	}

	return nil
}

func (s *ChatService) stateByUserAndRoom(
	userStringUUID string,
	chatStringUUID string,
) (*RoomState, *dto.User, error) {
	chatUUID, errChat := uuid.Parse(chatStringUUID)
	userUUID, errUser := uuid.Parse(userStringUUID)
	if errChat != nil || errUser != nil {
		return nil, nil, errors.New("Wrong UUID. ")
	}
	roomState := s.roomStateProvider.GetRoomStateByUUID(chatUUID)
	if roomState == nil {
		return nil, nil, errors.New("Room not found. ")
	}
	if roomState.author.user.UUID == userUUID {
		return roomState, roomState.author.user, nil
	}
	if roomState.guest != nil && roomState.guest.user.UUID == userUUID {
		return roomState, roomState.guest.user, nil
	}
	return nil, nil, errors.New("Wrong combination of user and room. ")
}

func NewChatService(provider *RoomStateProvider) *ChatService {
	return &ChatService{
		roomStateProvider: provider,
	}
}
