package services

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/internal/core/repo"
)

// RoomService lol, can I have a cheeseburger ?
type RoomService struct {
	managerProvider *RoomManagerProvider
	repo            *repo.Repository
	txLock          sync.RWMutex
}

func (r *RoomService) Create(user *domain.User) (*domain.Room, error) {
	// "Create" must be atomic transaction, using mutex
	r.txLock.Lock()
	defer r.txLock.Unlock()
	if r.repo.IsGuest(user) {
		r.leaveAsGuest(user)
	}
	room := r.repo.CreateRoomForUser(user)
	r.managerProvider.makeRoomManager(room)
	return room, nil
}

func (r *RoomService) Join(roomUuid uuid.UUID, user *domain.User) (*domain.Room, error) {
	r.txLock.Lock()
	defer r.txLock.Unlock()

	room := r.repo.FindRoomByUuid(roomUuid)
	if room == nil {
		return nil, errors.New("Room not exists. ")
	}
	if room.Guest == user {
		// already joined
		return room, nil
	}
	if room.Author == user {
		return nil, errors.New("User cant join as guest to his room. ")
	}
	if room.Guest != nil {
		return nil, errors.New("Room occupied. ")
	}

	manager := r.managerProvider.GetRoomManager(room)
	if manager == nil || !manager.IsAlive() {
		// room in state of deletion
		return nil, errors.New("Room in process of deletion. ")
	}

	r.repo.CommitUserJoin(room, user)
	manager.JoinRoom(user)
	return room, nil
}

func (r *RoomService) State(user *domain.User) *domain.Room {
	return r.repo.FindRoomByUser(user)
}

func (r *RoomService) Leave(user *domain.User) {
	r.txLock.Lock()
	defer r.txLock.Unlock()

	room := r.repo.FindRoomByUser(user)
	if room.Author == user {
		r.repo.CommitRoomShutdown(room)
		r.managerProvider.Close(room)
		return
	}
	r.leaveAsGuest(user)
}

func (r *RoomService) leaveAsGuest(user *domain.User) {
	room := r.repo.FindRoomByUser(user)
	r.repo.CommitUserLeave(room)
	manager := r.managerProvider.GetRoomManager(room)
	if manager != nil {
		manager.GuestLeave(user)
	}
}

func (r *RoomService) GetRoom(user *domain.User, stringUuid string) (*domain.Room, error) {
	roomUuid, err := uuid.Parse(stringUuid)
	if err != nil {
		return nil, err
	}
	room := r.State(user)
	if room != nil && room.UUID == roomUuid {
		return room, nil
	}
	return nil, errors.New("Wrong room. ")
}

func NewRoomService(provider *RoomManagerProvider, repo *repo.Repository) *RoomService {
	return &RoomService{
		managerProvider: provider,
		repo:            repo,
		txLock:          sync.RWMutex{},
	}
}
