package services

import (
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/internal/core/roommanager"
)

type RoomManagerProvider struct {
	mu    sync.RWMutex
	rooms map[uuid.UUID]*roommanager.RoomManager
}

func (r *RoomManagerProvider) makeRoomManager(room *domain.Room) {
	r.mu.Lock()
	manager := r.GetRoomManager(room)
	if manager == nil || !manager.IsAlive() {
		manager = roommanager.MakeManager(room.Author)
		r.rooms[room.UUID] = manager
	}
	r.mu.Unlock()
}

func (r *RoomManagerProvider) GetRoomManager(room *domain.Room) *roommanager.RoomManager {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if manager, exists := r.rooms[room.UUID]; exists {
		return manager
	}
	return nil
}

func (r *RoomManagerProvider) Close(room *domain.Room) {
	manager := r.GetRoomManager(room)
	r.mu.Lock()
	defer r.mu.Unlock()
	if manager != nil {
		manager.Shutdown()
		delete(r.rooms, room.UUID)
	}
}

func NewRoomProvider() *RoomManagerProvider {
	return &RoomManagerProvider{
		mu:    sync.RWMutex{},
		rooms: make(map[uuid.UUID]*roommanager.RoomManager),
	}
}
