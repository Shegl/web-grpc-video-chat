package inroom

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
	"web-grpc-video-chat/src/inroom/stream"
)

type RoomState struct {
	isAlive bool
	room    *dto.Room
	mu      sync.RWMutex
	author  *UserState
	guest   *UserState
	chat    *ChatState

	roomCtx context.Context
	Close   func()
}

type UserState struct {
	user         *dto.User
	state        *stream.User
	inputStream  *websocket.Conn
	outputStream *AVStream
	stateStream  *StateStream
	chatStream   *ChatStream
}

type ChatState struct {
	room     *dto.Room
	messages []*chat.ChatMessage
	msgChan  chan *chat.ChatMessage
	mu       sync.RWMutex
}

type ChatStream struct {
	stream  chat.Chat_ListenServer
	closeCh chan struct{}
}

type StateStream struct {
	stream  stream.Stream_StreamStateServer
	closeCh chan struct{}
}

type AVStream struct {
	stream  stream.Stream_AVStreamServer
	closeCh chan struct{}
}

func (roomState *RoomState) JoinRoomUpdate(guest *dto.User) error {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	if roomState.guest != nil && roomState.guest.user != guest {
		return errors.New("Cant update state, guest slot is occupied. ")
	}
	if roomState.guest != nil && roomState.guest.user == guest {
		return nil
	}
	roomState.guest = NewUserState(guest)
	return nil
}

func (roomState *RoomState) GetUserState(user *dto.User) *UserState {
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	return roomState.getUserState(user)
}

func (roomState *RoomState) getUserState(user *dto.User) *UserState {
	var userState *UserState
	if roomState.author.user == user {
		userState = roomState.author
	}
	if roomState.guest != nil && roomState.guest.user == user {
		userState = roomState.guest
	}
	return userState
}

func (roomState *RoomState) UpdateUserStreamState(user *dto.User, userStreamState *stream.User) {
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		userState.state = userStreamState
	}
}

func (roomState *RoomState) RoomChatConnect(
	user *dto.User,
	server chat.Chat_ListenServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		close(userState.chatStream.closeCh)
		userState.chatStream.stream = server
		userState.chatStream.closeCh = make(chan struct{})
		return userState.chatStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) AVStreamConnect(
	user *dto.User,
	server stream.Stream_AVStreamServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		close(userState.outputStream.closeCh)
		userState.outputStream.stream = server
		userState.outputStream.closeCh = make(chan struct{})
		return userState.outputStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) StateStreamConnect(
	user *dto.User,
	server stream.Stream_StreamStateServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		close(userState.stateStream.closeCh)
		userState.stateStream.stream = server
		userState.stateStream.closeCh = make(chan struct{})
		return userState.stateStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) LeaveRoomUpdate(guest *dto.User) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	roomState.closeChannels(guest)
}

func (roomState *RoomState) closeChannels(user *dto.User) {
	if roomState.guest != nil && roomState.guest.user == user {
		close(roomState.guest.chatStream.closeCh)
		close(roomState.guest.stateStream.closeCh)
		close(roomState.guest.outputStream.closeCh)
		roomState.guest = nil
	}
	if roomState.author.user == user {
		close(roomState.author.chatStream.closeCh)
		close(roomState.author.stateStream.closeCh)
		close(roomState.author.outputStream.closeCh)
		roomState.author = nil
	}
}

func NewUserState(user *dto.User) *UserState {
	return &UserState{
		user: user,
		state: &stream.User{
			IsCamEnabled: false,
			IsMuted:      true,
			UserUUID:     "",
			UserName:     "",
		},
		inputStream:  nil,
		outputStream: &AVStream{stream: nil, closeCh: make(chan struct{})},
		stateStream:  &StateStream{stream: nil, closeCh: make(chan struct{})},
		chatStream:   &ChatStream{stream: nil, closeCh: make(chan struct{})},
	}
}

func NewRoomState(room *dto.Room) *RoomState {
	roomCtx, cancelFunc := context.WithCancel(context.Background())
	roomState := &RoomState{
		room:    room,
		isAlive: true,
		mu:      sync.RWMutex{},
		author:  NewUserState(room.Author),
		guest:   nil,
		roomCtx: roomCtx,
		Close:   cancelFunc,
		chat:    nil,
	}
	AddChatState(roomState)
	return roomState
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
