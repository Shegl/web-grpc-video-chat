package roommanager

import (
	"sync"
	"web-grpc-video-chat/src/internal/core/domain"
)

type chatState struct {
	messages []domain.ChatMessage
	mu       sync.RWMutex
}

func (c *chatState) getMessages() []domain.ChatMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.messages[:len(c.messages)]
}

func (c *chatState) appendMessage(message domain.ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, message)
}

func makeChatState() *chatState {
	return &chatState{
		mu:       sync.RWMutex{},
		messages: []domain.ChatMessage{},
	}
}
