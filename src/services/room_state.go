package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"sync"
	"time"
	"web-grpc-video-chat/src/chat"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/streams"
)

type RoomStateProvider struct {
	mu         sync.RWMutex
	roomStates map[uuid.UUID]*RoomState
}

func (r *RoomStateProvider) MakeRoomState(room *dto.Room, chatState *ChatState) *RoomState {
	roomCtx, cancelFunc := context.WithCancel(context.Background())
	r.mu.Lock()
	roomState := &RoomState{
		room:    room,
		isAlive: true,
		mu:      sync.RWMutex{},
		author: &UserState{
			user:   room.Author,
			stream: nil,
			streamState: &streams.User{
				IsCamEnabled: false,
				IsMuted:      true,
				UserUUID:     "",
				UserName:     "",
			},
			stateServer:  nil,
			streamServer: nil,
			chatServer:   nil,
		},
		guest:   nil,
		roomCtx: roomCtx,
		Close:   cancelFunc,
		chat:    chatState,
	}
	r.roomStates[room.UUID] = roomState
	r.mu.Unlock()
	return roomState
}

func (r *RoomStateProvider) GetRoomStateByUUID(roomUUID uuid.UUID) *RoomState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if roomState, exists := r.roomStates[roomUUID]; exists {
		return roomState
	}
	return nil
}

func (r *RoomStateProvider) GetRoomState(room *dto.Room) *RoomState {
	return r.GetRoomStateByUUID(room.UUID)
}

func (r *RoomStateProvider) Forget(roomState *RoomState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	if !roomState.isAlive {
		return
	}
	roomState.isAlive = false
	roomState.Close()
	delete(r.roomStates, roomState.room.UUID)
	go func(state *RoomState) {
		time.Sleep(time.Second * 10)
		close(state.chat.msgChan)
	}(roomState)
}

func NewRoomStateProvider() *RoomStateProvider {
	return &RoomStateProvider{
		mu:         sync.RWMutex{},
		roomStates: make(map[uuid.UUID]*RoomState),
	}
}

type RoomState struct {
	isAlive bool
	room    *dto.Room
	mu      sync.RWMutex
	author  *UserState
	guest   *UserState

	chat *ChatState

	roomCtx context.Context
	Close   func()
}

type UserState struct {
	user         *dto.User
	stream       *websocket.Conn
	streamState  *streams.User
	stateServer  streams.Stream_StreamStateServer
	streamServer streams.Stream_AVStreamServer
	chatServer   chat.Chat_ListenServer
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
	roomState.guest = &UserState{
		user:   guest,
		stream: nil,
		streamState: &streams.User{
			IsCamEnabled: false,
			IsMuted:      true,
			UserUUID:     "",
			UserName:     "",
		},
		stateServer:  nil,
		streamServer: nil,
		chatServer:   nil,
	}
	return nil
}

func (roomState *RoomState) LeaveRoomUpdate(guest *dto.User) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	if roomState.guest.user == guest {
		roomState.guest = nil
	}
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

func (roomState *RoomState) UpdateUserStreamState(user *dto.User, userStreamState *streams.User) {
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		userState.streamState = userStreamState
	}
}

func (roomState *RoomState) RoomChatConnected(
	user *dto.User,
	server chat.Chat_ListenServer,
) error {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		userState.chatServer = server
		return nil
	}
	return errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) AVStreamConnected(
	user *dto.User,
	server streams.Stream_AVStreamServer,
) error {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		userState.streamServer = server
		return nil
	}
	return errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) StateStreamConnected(
	user *dto.User,
	server streams.Stream_StreamStateServer,
) error {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		userState.stateServer = server
		return nil
	}
	return errors.New("Cant update state, user not in room. ")
}
