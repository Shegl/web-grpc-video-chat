package services

import (
	"errors"
	"github.com/google/uuid"
	"macos-cam-grpc-chat/src/dto"
	"sync"
)

type RoomService struct {
	rooms    map[uuid.UUID]*dto.Room
	asGuests map[uuid.UUID]*dto.Room
	mu       sync.RWMutex
}

func (r *RoomService) Create(user *dto.User) (*dto.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.rooms[user.UUID]; exists {
		// room already exists
		// let user get back to his room,
		// if he needs new room, he will get that he needs leave his room first
		return room, nil
	}
	if _, exists := r.asGuests[user.UUID]; exists {
		// let's check user, what if he already as guest in someone else room
		r.leave(user)
	}
	return r.create(user), nil
}

func (r *RoomService) create(user *dto.User) *dto.Room {
	room := &dto.Room{
		State:  1,
		UUID:   uuid.New(),
		Author: user,
		Guest:  nil,
	}
	r.rooms[user.UUID] = room
	return room
}

func (r *RoomService) Join(roomUUID uuid.UUID, user *dto.User) (*dto.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.rooms[user.UUID]; exists {
		// room exists, good. Is it occupied?
		if room.Guest == nil {
			return r.join(roomUUID, user), nil
		} else if room.Guest == user {
			// same user, its counter as rejoin
			return room, nil
		} else {
			return nil, errors.New("Room occupied. ")
		}
	}
	return nil, errors.New("Room not exists. ")
}

func (r *RoomService) join(roomUUID uuid.UUID, user *dto.User) *dto.Room {
	if room, exists := r.rooms[user.UUID]; exists {
		room.Guest = user
		return room
	}
	return nil
}

func (r *RoomService) Leave(user *dto.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.asGuests[user.UUID]; exists {
		room.Guest = nil
		delete(r.asGuests, user.UUID)
	}
}

func (r *RoomService) leave(user *dto.User) {
	if room, exists := r.asGuests[user.UUID]; exists {
		room.Guest = nil
		delete(r.asGuests, user.UUID)
	}
}

func NewRoomService() *RoomService {
	return &RoomService{
		rooms:    make(map[uuid.UUID]*dto.Room),
		asGuests: make(map[uuid.UUID]*dto.Room),
		mu:       sync.RWMutex{},
	}
}
