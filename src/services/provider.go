package services

import (
	"context"
	"github.com/google/uuid"
	"sync"
	"time"
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
			user: room.Author,
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
		// close channels
		// for guest
		if state.guest != nil && state.guest.chatStream.stream != nil {
			state.guest.chatStream.stream = nil
			close(state.guest.chatStream.closeCh)
		}
		if state.guest != nil && state.guest.stateStream.stream != nil {
			state.guest.stateStream.stream = nil
			close(state.guest.stateStream.closeCh)
		}
		if state.guest != nil && state.guest.outputStream.stream != nil {
			state.guest.outputStream.stream = nil
			close(state.guest.outputStream.closeCh)
		}
		// for author
		if state.author.chatStream.stream != nil {
			state.author.chatStream.stream = nil
			close(state.guest.chatStream.closeCh)
		}
		if state.author.stateStream.stream != nil {
			state.author.stateStream.stream = nil
			close(state.guest.stateStream.closeCh)
		}
		if state.author.outputStream.stream != nil {
			state.author.outputStream.stream = nil
			close(state.guest.outputStream.closeCh)
		}
		close(state.chat.msgChan)
	}(roomState)
}

func NewRoomStateProvider() *RoomStateProvider {
	return &RoomStateProvider{
		mu:         sync.RWMutex{},
		roomStates: make(map[uuid.UUID]*RoomState),
	}
}
