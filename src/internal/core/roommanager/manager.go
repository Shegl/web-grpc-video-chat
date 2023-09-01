package roommanager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/pb/stream"
)

type managerEvent int

const (
	eventRoomCreated  managerEvent = 0
	eventRoomShutdown              = 1
	eventRoomJoin                  = 2
	eventRoomLeave                 = 3
)

type RoomManager struct {
	isAlive bool
	mu      sync.RWMutex
	author  *userSlot
	guest   *userSlot
	chat    *chatState
	roomCtx context.Context
	close   func()

	observers []ManagerObserver
}

func (m *RoomManager) open() {
	if m.isAlive {
		return
	}
	m.isAlive = true
	go m.notifyObservers(eventRoomCreated, *m.author.user)
}

func (m *RoomManager) IsAlive() bool {
	return m.isAlive
}

func (m *RoomManager) JoinRoom(guest *domain.User) {
	m.mu.Lock()
	m.guest.takeSlot(guest)
	m.mu.Unlock()
	go m.notifyObservers(eventRoomJoin, *guest)
}

func (m *RoomManager) GuestLeave(guest *domain.User) {
	m.mu.Lock()
	if m.guest.isUser(guest) {
		m.guest.freeUpSlot()
		go m.notifyObservers(eventRoomLeave, *guest)
	}
	m.mu.Unlock()
}

func (m *RoomManager) InRoom(user *domain.User) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.guest.isUser(user) || m.author.isUser(user)
}

func (m *RoomManager) getUserSlot(user *domain.User) *userSlot {
	var slot *userSlot
	if m.author.isUser(user) {
		slot = m.author
	}
	if m.guest != nil && m.author.isUser(user) {
		slot = m.guest
	}
	return slot
}

func (m *RoomManager) avStreamConnect(
	user *domain.User,
	server stream.Stream_AVStreamServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.setAVStream(server), nil
	}
	// very rare but possible error must be handled
	return nil, errors.New("Cant update state, user already left room nanoseconds ago. ")
}

func (m *RoomManager) stateStreamConnect(
	user *domain.User,
	server stream.Stream_StreamStateServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.setStateStream(server), nil
	}
	// very rare but possible error must be handled
	return nil, errors.New("Cant update state, user already left room nanoseconds ago. ")
}

func (m *RoomManager) updateUserState(user *domain.User, userStreamState *stream.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		slot.updateUserState(userStreamState)
	}
}

func (m *RoomManager) RoomChatConnect(
	user *domain.User,
	server ChatListenerServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.setChatStream(server), nil
	}
	// very rare but possible error must be handled
	return nil, errors.New("Cant update state, user already left room nanoseconds ago. ")
}

func (m *RoomManager) GetChatHistory() []domain.ChatMessage {
	return m.chat.getMessages()
}

func (m *RoomManager) ChatBroadcast(message domain.ChatMessage) {
	m.chat.appendMessage(message)
	m.mu.RLock()
	m.author.sendChatMessage(message)
	m.guest.sendChatMessage(message)
	m.mu.RUnlock()
}

func (m *RoomManager) RegisterObserver(observer ManagerObserver) {
	m.observers = append(m.observers, observer)
}

func (m *RoomManager) notifyObservers(event managerEvent, user domain.User) {
	for _, observer := range m.observers {
		switch event {
		case eventRoomCreated:
			observer.RoomCreate(m, user)
			break
		case eventRoomShutdown:
			observer.RoomShutdown(m, user)
			break
		case eventRoomJoin:
			observer.RoomUserJoin(m, user)
			break
		case eventRoomLeave:
			observer.RoomUserLeave(m, user)
			break
		default:
			panic(errors.New(fmt.Sprintf("Unhandled event in room manager, are you stupid?: %v", event)))
		}
	}
}

func (m *RoomManager) Shutdown() {
	if !m.isAlive {
		return
	}
	m.mu.Lock()
	go m.notifyObservers(eventRoomShutdown, *m.author.user)
	defer m.mu.Unlock()
	m.isAlive = false
	m.close()
	m.guest.freeUpSlot()
	go func() {
		time.Sleep(time.Second * 5)
		m.author.freeUpSlot()
		m.observers = nil
	}()
}
