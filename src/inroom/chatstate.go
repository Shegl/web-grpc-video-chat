package inroom

import (
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
)

type chatState struct {
	room     *dto.Room
	messages []*chat.ChatMessage
	mu       sync.RWMutex
}

func (c *chatState) getMessages() []*chat.ChatMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	messagesCopy := make([]*chat.ChatMessage, len(c.messages))
	copy(messagesCopy, c.messages)
	return messagesCopy
}

func (c *chatState) appendMessage(message *chat.ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, message)
}
