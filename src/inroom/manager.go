package inroom

import (
	"context"
	"errors"
	"sync"
	"time"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
	"web-grpc-video-chat/src/inroom/stream"
)

type RoomManager struct {
	isAlive bool
	room    *dto.Room
	mu      sync.RWMutex
	author  *userSlot
	guest   *userSlot
	chat    *chatState
	roomCtx context.Context
	close   func()
}

func (m *RoomManager) JoinRoom(guest *dto.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.guest.takeSlot(guest)
}

func (m *RoomManager) GuestLeave(guest *dto.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.guest.isUser(guest) {
		m.guest.freeUpSlot()
	}
}

func (m *RoomManager) IsAlive() bool {
	return m.isAlive
}

func (m *RoomManager) inRoom(user *dto.User) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.guest.isUser(user) || m.author.isUser(user)
}

func (m *RoomManager) updateUserState(user *dto.User, userStreamState *stream.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		slot.updateUserState(userStreamState)
	}
}

func (m *RoomManager) roomChatConnect(
	user *dto.User,
	server chat.Chat_ListenServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.connectChatStream(server), nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (m *RoomManager) getChatHistory() []*chat.ChatMessage {
	return m.chat.getMessages()
}

func (m *RoomManager) chatBroadcast(message *chat.ChatMessage) {
	m.chat.appendMessage(message)
	m.mu.RLock()
	m.author.sendChatMessage(message)
	m.guest.sendChatMessage(message)
	m.mu.RUnlock()
}

func (m *RoomManager) avStreamConnect(
	user *dto.User,
	server stream.Stream_AVStreamServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.connectAVStream(server), nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (m *RoomManager) stateStreamConnect(
	user *dto.User,
	server stream.Stream_StreamStateServer,
) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	slot := m.getUserSlot(user)
	if slot != nil {
		return slot.connectStateStream(server), nil
	}
	return nil, errors.New("Cant update state, user not in room. ")
}

func (m *RoomManager) getUserSlot(user *dto.User) *userSlot {
	var slot *userSlot
	if m.author.isUser(user) {
		slot = m.author
	}
	if m.guest != nil && m.author.isUser(user) {
		slot = m.guest
	}
	return slot
}

func (m *RoomManager) shutdown() {
	if !m.isAlive {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isAlive = false
	m.close()
	m.guest.freeUpSlot()
	go func() {
		time.Sleep(time.Second * 5)
		m.author.freeUpSlot()
	}()
}
