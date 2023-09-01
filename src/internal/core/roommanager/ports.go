package roommanager

import "web-grpc-video-chat/src/internal/core/domain"

type ChatListenerServer interface {
	Send(msg domain.ChatMessage) error
}

type ManagerObserver interface {
	RoomCreate(manager *RoomManager, user domain.User)
	RoomShutdown(manager *RoomManager, user domain.User)
	RoomUserJoin(manager *RoomManager, user domain.User)
	RoomUserLeave(manager *RoomManager, user domain.User)
}
