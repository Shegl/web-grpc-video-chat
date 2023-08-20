package dto

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type lookupStorage struct {
	// rooms
	rooms    map[uuid.UUID]*Room
	asAuthor map[uuid.UUID]*Room
	asGuest  map[uuid.UUID]*Room

	// users
	authUsers map[uuid.UUID]*User
}

type Repository struct {
	mu sync.RWMutex
	ls lookupStorage
}

func (r *Repository) CreateUser(userName string) (*User, error) {
	userUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	if _, exists := r.ls.authUsers[userUuid]; exists {
		r.mu.Unlock()
		return nil, errors.New("UUID taken. ")
	}
	user := &User{
		Name: userName,
		UUID: userUuid,
	}
	r.ls.authUsers[userUuid] = user
	r.mu.Unlock()
	return user, nil
}

func (r *Repository) CreateRoomForUser(user *User) *Room {
	roomUuid, err := uuid.NewRandom()
	if err != nil {
		return r.CreateRoomForUser(user)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.ls.asAuthor[user.UUID]; exists {
		return room
	}
	room := &Room{
		UUID:   roomUuid,
		Author: user,
		Guest:  nil,
	}
	r.ls.rooms[room.UUID] = room
	r.ls.asAuthor[user.UUID] = room
	return room
}

func (r *Repository) FindRoomByUuid(roomUuid uuid.UUID) *Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ls.rooms[roomUuid]
}

func (r *Repository) FindRoomByUser(user *User) *Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if room, exists := r.ls.asAuthor[user.UUID]; exists {
		return room
	}
	if room, exists := r.ls.asGuest[user.UUID]; exists {
		return room
	}
	return nil
}

func (r *Repository) FindRoomByString(stringUUID string) (*Room, error) {
	roomUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	room := r.FindRoomByUuid(roomUUID)
	return room, nil
}

func (r *Repository) CommitUserJoin(room *Room, user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room.Guest = user
	r.ls.asGuest[user.UUID] = room
}

func (r *Repository) CommitGuestLeave(room *Room) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.ls.asGuest, room.Guest.UUID)
	room.Guest = nil
}

func (r *Repository) CommitRoomShutdown(room *Room) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if room.Guest != nil {
		delete(r.ls.asGuest, room.Guest.UUID)
	}
	delete(r.ls.asAuthor, room.Author.UUID)
	delete(r.ls.rooms, room.UUID)
	room.Guest = nil
	room.Author = nil
}

func (r *Repository) IsAuthor(user *User) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.ls.asAuthor[user.UUID]
	return exists
}

func (r *Repository) IsGuest(user *User) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.ls.asGuest[user.UUID]
	return exists
}

func (r *Repository) FindUserByUuid(uuid uuid.UUID) *User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ls.authUsers[uuid]
}

func (r *Repository) FindUserByString(stringUUID string) (*User, error) {
	userUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	return r.FindUserByUuid(userUUID), nil
}

func (r *Repository) ForgetUserByUuid(userUUID uuid.UUID) {
	r.mu.Lock()
	delete(r.ls.authUsers, userUUID)
	r.mu.Unlock()
}

func NewRepository() *Repository {
	return &Repository{
		mu: sync.RWMutex{},
		ls: lookupStorage{
			rooms:     make(map[uuid.UUID]*Room),
			asAuthor:  make(map[uuid.UUID]*Room),
			asGuest:   make(map[uuid.UUID]*Room),
			authUsers: make(map[uuid.UUID]*User),
		}}
}
