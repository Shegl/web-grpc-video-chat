package services

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom"
)

type RoomService struct {
	roomProvider *inroom.RoomProvider
	repo         *dto.Repository
	txLock       sync.RWMutex
}

func (r *RoomService) Create(user *dto.User) (*dto.Room, error) {
	// "Create" must be atomic transaction, using mutex
	r.txLock.Lock()
	defer r.txLock.Unlock()
	if r.repo.IsGuest(user) {
		r.leaveAsGuest(user)
	}
	return r.create(user), nil
}

func (r *RoomService) create(user *dto.User) *dto.Room {
	room := r.repo.CreateRoomForUser(user)
	r.roomProvider.MakeRoom(room)
	return room
}

func (r *RoomService) Join(roomUuid uuid.UUID, user *dto.User) (*dto.Room, error) {
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
	return r.join(room, user), nil
}

func (r *RoomService) join(room *dto.Room, user *dto.User) *dto.Room {
	manager := r.roomProvider.GetRoomManager(room)
	if manager == nil || !manager.IsAlive() {
		// room in state of deletion
		return nil
	}
	err := manager.JoinRoom(user)
	if err != nil {
		panic(err)
	}
	r.repo.CommitUserJoin(room, user)
	return room
}

func (r *RoomService) State(user *dto.User) *dto.Room {
	return r.repo.FindRoomByUser(user)
}

func (r *RoomService) Leave(user *dto.User) {
	r.txLock.Lock()
	defer r.txLock.Unlock()

	room := r.repo.FindRoomByUser(user)
	if room.Author == user {
		r.roomProvider.Close(room)
		r.repo.CommitRoomShutdown(room)
		return
	}
	r.leaveAsGuest(user)
}

func (r *RoomService) leaveAsGuest(user *dto.User) {
	room := r.repo.FindRoomByUser(user)
	if room == nil {
		return
	}
	err := r.repo.CommitUserLeave(room, user)
	if err != nil {
		panic(err)
	}
	manager := r.roomProvider.GetRoomManager(room)
	if manager != nil {
		manager.GuestLeave(user)
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

func NewRoomService(provider *inroom.RoomProvider, repo *dto.Repository) *RoomService {
	return &RoomService{
		roomProvider: provider,
		repo:         repo,
		txLock:       sync.RWMutex{},
	}
}
