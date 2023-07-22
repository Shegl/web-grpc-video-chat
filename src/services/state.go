package services

import (
	"context"
	"errors"
	"golang.org/x/net/websocket"
	"sync"
	"web-grpc-video-chat/src/chat"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/streams"
)

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
	state        *streams.User
	inputStream  *websocket.Conn
	outputStream *AVStream
	stateStream  *StateStream
	chatStream   *ChatStream
}

type ChatStream struct {
	stream  chat.Chat_ListenServer
	closeCh chan struct{}
}

type StateStream struct {
	stream  streams.Stream_StreamStateServer
	closeCh chan struct{}
}

type AVStream struct {
	stream  streams.Stream_AVStreamServer
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
	roomState.guest = &UserState{
		user: guest,
		state: &streams.User{
			IsCamEnabled: false,
			IsMuted:      true,
			UserUUID:     "",
			UserName:     "",
		},
		inputStream:  nil,
		outputStream: &AVStream{},
		stateStream:  &StateStream{},
		chatStream:   &ChatStream{},
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
		userState.state = userStreamState
	}
}

func (roomState *RoomState) RoomChatConnected(
	user *dto.User,
	server chat.Chat_ListenServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		if userState.chatStream.stream != nil {
			userState.chatStream.stream = nil
			close(userState.stateStream.closeCh)
		}
		userState.chatStream.stream = server
		userState.chatStream.closeCh = make(chan struct{})
		return userState.chatStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) AVStreamConnected(
	user *dto.User,
	server streams.Stream_AVStreamServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		if userState.outputStream.stream != nil {
			userState.outputStream.stream = nil
			close(userState.outputStream.closeCh)
		}
		userState.outputStream.stream = server
		userState.outputStream.closeCh = make(chan struct{})
		return userState.outputStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (roomState *RoomState) StateStreamConnected(
	user *dto.User,
	server streams.Stream_StreamStateServer,
) (<-chan struct{}, error) {
	roomState.mu.Lock()
	defer roomState.mu.Unlock()
	userState := roomState.getUserState(user)
	if userState != nil {
		if userState.stateStream.stream != nil {
			userState.stateStream.stream = nil
			close(userState.stateStream.closeCh)
		}
		userState.stateStream.stream = server
		userState.stateStream.closeCh = make(chan struct{})
		return userState.stateStream.closeCh, nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}
