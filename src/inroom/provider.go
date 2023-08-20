package inroom

import (
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/dto"
)

type RoomProvider struct {
	mu    sync.RWMutex
	rooms map[uuid.UUID]*RoomManager
}

func (r *RoomProvider) MakeRoomManager(room *dto.Room) {
	r.mu.Lock()
	managedRoom := makeManager(room)
	r.rooms[room.UUID] = managedRoom
	r.mu.Unlock()
}

func (r *RoomProvider) GetRoomManager(room *dto.Room) *RoomManager {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if manager, exists := r.rooms[room.UUID]; exists {
		return manager
	}
	return nil
}

func (r *RoomProvider) Close(room *dto.Room) {
	manager := r.GetRoomManager(room)
	r.mu.Lock()
	defer r.mu.Unlock()
	if manager != nil {
		manager.shutdown()
		delete(r.rooms, room.UUID)
	}
}

func NewRoomProvider() *RoomProvider {
	return &RoomProvider{
		mu:    sync.RWMutex{},
		rooms: make(map[uuid.UUID]*RoomManager),
	}
}
