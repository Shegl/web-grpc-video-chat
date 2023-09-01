package roommanager

import (
	"context"
	"sync"
	"web-grpc-video-chat/src/internal/core/domain"
)

func MakeManager(author *domain.User) *RoomManager {
	roomCtx, cancelFunc := context.WithCancel(context.Background())
	manager := &RoomManager{
		mu:      sync.RWMutex{},
		roomCtx: roomCtx,
		close:   cancelFunc,
		author:  makeSlot(author),
		guest:   makeSlot(nil),
		chat:    makeChatState(),
	}
	manager.open()
	return manager
}
