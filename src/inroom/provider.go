package inroom

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
)

type RoomStateProvider struct {
	mu         sync.RWMutex
	roomStates map[uuid.UUID]*RoomState
}

func (r *RoomStateProvider) MakeRoomState(room *dto.Room) *RoomState {
	r.mu.Lock()
	roomState := NewRoomState(room)
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

func (r *RoomStateProvider) GetByUserAndRoom(
	userStringUUID string,
	roomStringUUID string,
) (*RoomState, *dto.User, error) {
	roomUUID, errChat := uuid.Parse(roomStringUUID)
	userUUID, errUser := uuid.Parse(userStringUUID)
	if errChat != nil || errUser != nil {
		return nil, nil, errors.New("Wrong UUID. ")
	}
	roomState := r.GetRoomStateByUUID(roomUUID)
	if roomState == nil {
		return nil, nil, errors.New("Room not found. ")
	}
	roomState.mu.RLock()
	defer roomState.mu.RUnlock()
	if roomState.author.user.UUID == userUUID {
		return roomState, roomState.author.user, nil
	}
	if roomState.guest != nil && roomState.guest.user.UUID == userUUID {
		return roomState, roomState.guest.user, nil
	}
	return nil, nil, errors.New("Wrong combination of user and room. ")
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
		// for author
		close(state.author.chatStream.closeCh)
		close(state.author.stateStream.closeCh)
		close(state.author.outputStream.closeCh)
		close(state.chat.msgChan)
	}(roomState)
}

func NewRoomStateProvider() *RoomStateProvider {
	return &RoomStateProvider{
		mu:         sync.RWMutex{},
		roomStates: make(map[uuid.UUID]*RoomState),
	}
}
