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

type UserState struct {
	user         *dto.User
	stream       *websocket.Conn
	streamState  *streams.User
	stateServer  streams.Stream_StreamStateServer
	streamServer streams.Stream_AVStreamServer
	chatServer   chat.Chat_ListenServer
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

type RoomStateProvider struct {
	mu         sync.RWMutex
	roomStates map[uuid.UUID]*RoomState
}

func (r *RoomStateProvider) MakeRoomState(room *dto.Room, author *dto.User, chatState *ChatState) *RoomState {
	roomCtx, cancelFunc := context.WithCancel(context.Background())
	r.mu.Lock()
	roomState := &RoomState{
		room:    room,
		isAlive: true,
		mu:      sync.RWMutex{},
		author: &UserState{
			user:   author,
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

func (r *RoomStateProvider) JoinRoomStateUpdate(room *dto.Room, guest *dto.User) (*RoomState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	roomState, exists := r.roomStates[room.UUID]
	if !exists {
		return nil, errors.New("Cant update state, state for this room not exists. ")
	}
	if roomState.guest != nil && roomState.guest.user != guest {
		return nil, errors.New("Cant update state, guest slot is occupied. ")
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
	return roomState, nil
}

func (r *RoomStateProvider) LeaveRoomStateUpdate(room *dto.Room, guest *dto.User) (*RoomState, error) {
	r.mu.Lock()
	roomState, exists := r.roomStates[room.UUID]
	if !exists {
		return nil, errors.New("Cant update state, state for this room not exists. ")
	}
	if roomState.guest.user == guest {
		roomState.guest = nil
	}
	r.mu.Unlock()
	return roomState, nil
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

func (r *RoomStateProvider) RoomChatConnected(
	room *dto.Room,
	user *dto.User,
	server chat.Chat_ListenServer,
) (*RoomState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userState, roomState := r.GetUserState(room, user)
	if userState != nil {
		userState.chatServer = server
		return roomState, nil
	}
	return nil, errors.New("Cant update state, no such room or user not in room. ")
}

func (r *RoomStateProvider) AVStreamConnected(
	room *dto.Room,
	user *dto.User,
	server streams.Stream_AVStreamServer,
) (*RoomState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userState, roomState := r.GetUserState(room, user)
	if userState != nil {
		userState.streamServer = server
		return roomState, nil
	}
	return nil, errors.New("Cant update state, no such room or user not in room. ")
}

func (r *RoomStateProvider) StateStreamConnected(
	room *dto.Room,
	user *dto.User,
	server streams.Stream_StreamStateServer,
) (*RoomState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userState, roomState := r.GetUserState(room, user)
	if userState != nil {
		userState.stateServer = server
		return roomState, nil
	}
	return nil, errors.New("Cant update state, no such room or user not in room. ")
}

func (r *RoomStateProvider) GetUserState(
	room *dto.Room,
	user *dto.User,
) (*UserState, *RoomState) {
	roomState := r.GetRoomState(room)
	if roomState != nil {
		roomState.mu.RLock()
		defer roomState.mu.RUnlock()
		var userState *UserState
		if room.Author == user {
			userState = roomState.author
		}
		if room.Guest == user {
			userState = roomState.guest
		}
		return userState, roomState
	}
	return nil, nil
}

func (r *RoomStateProvider) Forget(roomState *RoomState) {
	if !roomState.isAlive {
		return
	}
	r.mu.Lock()
	roomState.isAlive = false
	roomState.Close()
	delete(r.roomStates, roomState.room.UUID)
	r.mu.Unlock()
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
