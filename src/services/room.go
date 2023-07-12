package services

import (
	"errors"
	"github.com/google/uuid"
	"macos-cam-grpc-chat/src/dto"
	"sync"
)

type RoomService struct {
	rooms    map[uuid.UUID]*dto.Room
	asAuthor map[uuid.UUID]*dto.Room
	asGuest  map[uuid.UUID]*dto.Room
	mu       sync.RWMutex
}

func (r *RoomService) Create(user *dto.User) (*dto.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.asAuthor[user.UUID]; exists {
		// room already exists
		// let user get back to his room,
		// if he needs new room, he will get that he needs leave his room first
		return room, nil
	}
	if _, exists := r.asGuest[user.UUID]; exists {
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
	r.asAuthor[user.UUID] = room
	return room
}

func (r *RoomService) Join(roomUUID uuid.UUID, user *dto.User) (*dto.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.rooms[roomUUID]; exists {
		// room exists, good. Is it occupied?
		if room.Guest == nil && room.Author != user {
			return r.join(roomUUID, user), nil
		} else if room.Guest == user {
			// same user, its counter as rejoin
			return room, nil
		} else if room.Author == user {
			return nil, errors.New("User cant join as guest to his room. ")
		} else {
			return nil, errors.New("Room occupied. ")
		}
	}
	return nil, errors.New("Room not exists. ")
}

func (r *RoomService) join(roomUUID uuid.UUID, user *dto.User) *dto.Room {
	if room, exists := r.rooms[roomUUID]; exists {
		room.Guest = user
		return room
	}
	return nil
}

func (r *RoomService) State(user *dto.User) *dto.Room {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.asAuthor[user.UUID]; exists {
		return room
	}
	if room, exists := r.asGuest[user.UUID]; exists {
		return room
	}
	return nil
}

func (r *RoomService) Leave(user *dto.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.asAuthor[user.UUID]; exists {
		delete(r.rooms, room.UUID)
		delete(r.asAuthor, user.UUID)
		if room.Guest != nil {
			delete(r.asGuest, room.Guest.UUID)
		}
		return
	}
	if room, exists := r.asGuest[user.UUID]; exists {
		room.Guest = nil
		delete(r.asGuest, user.UUID)
	}
}

func (r *RoomService) leave(user *dto.User) {
	if room, exists := r.asGuest[user.UUID]; exists {
		room.Guest = nil
		delete(r.asGuest, user.UUID)
	}
}

func NewRoomService() *RoomService {
	return &RoomService{
		rooms:    make(map[uuid.UUID]*dto.Room),
		asAuthor: make(map[uuid.UUID]*dto.Room),
		asGuest:  make(map[uuid.UUID]*dto.Room),
		mu:       sync.RWMutex{},
	}
}
