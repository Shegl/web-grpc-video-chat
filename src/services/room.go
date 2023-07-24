package services

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom"
)

type RoomService struct {
	stateProvider *inroom.RoomStateProvider

	rooms    map[uuid.UUID]*dto.Room
	asAuthor map[uuid.UUID]*dto.Room
	asGuest  map[uuid.UUID]*dto.Room

	mu sync.RWMutex
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
	if _, existsAsGuest := r.asGuest[user.UUID]; existsAsGuest {
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
	r.rooms[room.UUID] = room
	r.asAuthor[user.UUID] = room

	r.stateProvider.MakeRoomState(room)

	return room
}

func (r *RoomService) Join(roomUUID uuid.UUID, user *dto.User) (*dto.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.rooms[roomUUID]; exists {
		// room exists, good. Is it occupied?
		if room.Guest == nil && room.Author != user {
			return r.join(roomUUID, user), nil
		}
		if room.Guest == user {
			// same user, its counter as rejoin
			return room, nil
		}
		// error handling
		var err error
		if room.Author == user {
			err = errors.New("User cant join as guest to his room. ")
		} else {
			err = errors.New("Room occupied. ")
		}
		return nil, err
	}
	return nil, errors.New("Room not exists. ")
}

func (r *RoomService) join(roomUUID uuid.UUID, user *dto.User) *dto.Room {
	if room, exists := r.rooms[roomUUID]; exists {
		roomState := r.stateProvider.GetRoomState(room)
		if roomState == nil {
			// room in state of deletion
			return nil
		}
		err := roomState.JoinRoomUpdate(user)
		if err != nil {
			panic(err)
		}
		room.Guest = user
		r.asGuest[user.UUID] = room
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
		roomState := r.stateProvider.GetRoomState(room)
		if roomState != nil {
			r.stateProvider.Forget(roomState)
		}
		delete(r.rooms, room.UUID)
		delete(r.asAuthor, user.UUID)
		return
	}
	r.leave(user)
}

func (r *RoomService) leave(user *dto.User) {
	if room, exists := r.asGuest[user.UUID]; exists {
		roomState := r.stateProvider.GetRoomState(room)
		if roomState != nil {
			roomState.LeaveRoomUpdate(user)
		}
		room.Guest = nil
		delete(r.asGuest, user.UUID)
	}
}

func (r *RoomService) GetRoom(user *dto.User, stringUUID string) (*dto.Room, error) {
	roomUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	room := r.State(user)
	if room != nil && room.UUID == roomUUID {
		return room, nil
	}
	return nil, errors.New("Wrong room. ")
}

func NewRoomService(provider *inroom.RoomStateProvider) *RoomService {
	return &RoomService{
		stateProvider: provider,

		rooms:    make(map[uuid.UUID]*dto.Room),
		asAuthor: make(map[uuid.UUID]*dto.Room),
		asGuest:  make(map[uuid.UUID]*dto.Room),
	}
}
