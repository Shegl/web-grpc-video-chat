package repo

import (
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/internal/core/domain"
)

type lookupStorage struct {
	// rooms
	rooms    map[uuid.UUID]*domain.Room
	asAuthor map[uuid.UUID]*domain.Room
	asGuest  map[uuid.UUID]*domain.Room

	// users
	authUsers map[uuid.UUID]*domain.User
}

type Repository struct {
	mu sync.RWMutex
	ls lookupStorage
}

func (r *Repository) CreateUser(userName string) (*domain.User, error) {
	userUuid, err := uuid.NewRandom()
	if err != nil {
		return r.CreateUser(userName)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	user := &domain.User{
		Name: userName,
		UUID: userUuid,
	}
	r.ls.authUsers[userUuid] = user
	return user, nil
}

func (r *Repository) CreateRoomForUser(user *domain.User) *domain.Room {
	roomUuid, err := uuid.NewRandom()
	if err != nil {
		return r.CreateRoomForUser(user)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if room, exists := r.ls.asAuthor[user.UUID]; exists {
		return room
	}
	room := &domain.Room{
		UUID:   roomUuid,
		Author: user,
		Guest:  nil,
	}
	r.ls.rooms[room.UUID] = room
	r.ls.asAuthor[user.UUID] = room
	return room
}

func (r *Repository) FindRoomByUuid(roomUuid uuid.UUID) *domain.Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ls.rooms[roomUuid]
}

func (r *Repository) FindRoomByUser(user *domain.User) *domain.Room {
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

func (r *Repository) FindRoomByString(stringUUID string) (*domain.Room, error) {
	roomUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	room := r.FindRoomByUuid(roomUUID)
	return room, nil
}

func (r *Repository) CommitUserJoin(room *domain.Room, user *domain.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room.Guest = user
	r.ls.asGuest[user.UUID] = room
}

func (r *Repository) CommitUserLeave(room *domain.Room) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.ls.asGuest, room.Guest.UUID)
	room.Guest = nil
}

func (r *Repository) CommitRoomShutdown(room *domain.Room) {
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

func (r *Repository) IsAuthor(user *domain.User) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.ls.asAuthor[user.UUID]
	return exists
}

func (r *Repository) IsGuest(user *domain.User) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.ls.asGuest[user.UUID]
	return exists
}

func (r *Repository) FindUserByUuid(uuid uuid.UUID) *domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ls.authUsers[uuid]
}

func (r *Repository) FindUserByString(stringUUID string) (*domain.User, error) {
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
			rooms:     make(map[uuid.UUID]*domain.Room),
			asAuthor:  make(map[uuid.UUID]*domain.Room),
			asGuest:   make(map[uuid.UUID]*domain.Room),
			authUsers: make(map[uuid.UUID]*domain.User),
		}}
}
