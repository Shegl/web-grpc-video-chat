package streams

import (
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
)

type StreamStateProvider struct {
	mu           sync.RWMutex
	streamStates map[uuid.UUID]*RoomStreamState
}

type RoomStreamState struct {
	room *dto.Room
	mu   sync.RWMutex

	author *UserStreamState
	guest  *UserStreamState
}

type UserStreamState struct {
	state        *User
	stream       *websocket.Conn
	stateServer  Stream_StreamStateServer
	streamServer Stream_AVStreamServer
}

func (s *StreamStateProvider) MakeRoomStreamState(room *dto.Room) *RoomStreamState {
	roomState := &RoomStreamState{
		room: room,
		mu:   sync.RWMutex{},
		author: &UserStreamState{
			state:        nil,
			stream:       nil,
			stateServer:  nil,
			streamServer: nil,
		},
		guest: &UserStreamState{
			state:        nil,
			stream:       nil,
			stateServer:  nil,
			streamServer: nil,
		},
	}
	s.streamStates[room.UUID] = roomState
	return roomState
}

func (s *StreamStateProvider) GetRoomState(room *dto.Room) *RoomStreamState {
	s.mu.Lock()
	roomState, exists := s.streamStates[room.UUID]
	if !exists {
		roomState = s.MakeRoomStreamState(room)
	}
	s.mu.Unlock()
	return roomState
}

func (s *StreamStateProvider) GetUserStreamState(user *dto.User, room *dto.Room) (*UserStreamState, *RoomStreamState) {
	roomState := s.GetRoomState(room)
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	var userState *UserStreamState
	if room.Author == user {
		userState = roomState.author
	} else {
		userState = roomState.guest
	}
	return userState, roomState
}

func (s *StreamStateProvider) SendStateUpdates(roomState *RoomStreamState) {
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	stateMessage := &StateMessage{
		Time:   time.Now().Unix(),
		UUID:   uuid.NewString(),
		Author: roomState.author.state,
		Guest:  roomState.guest.state,
	}
	if roomState.author.stateServer != nil {
		roomState.author.stateServer.Send(stateMessage)
	}
	if roomState.guest.stateServer != nil {
		roomState.guest.stateServer.Send(stateMessage)
	}
}

func MakeStreamStateProvider() *StreamStateProvider {
	return &StreamStateProvider{
		mu:           sync.RWMutex{},
		streamStates: make(map[uuid.UUID]*RoomStreamState),
	}
}
