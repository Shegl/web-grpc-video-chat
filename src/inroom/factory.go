package inroom

import (
	"context"
	"github.com/google/uuid"
	"sync"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
	"web-grpc-video-chat/src/inroom/stream"
)

func makeManager(room *dto.Room) *RoomManager {
	roomCtx, cancelFunc := context.WithCancel(context.Background())
	manager := &RoomManager{
		room:    room,
		isAlive: true,
		mu:      sync.RWMutex{},
		roomCtx: roomCtx,
		close:   cancelFunc,
		chat:    nil,
	}
	authorSlot, guestSlot := makeSlots(room.Author, manager)
	manager.author = authorSlot
	manager.guest = guestSlot
	manager.chat = makeChatState()
	return manager
}

func makeSlots(author *dto.User, manager *RoomManager) (*userSlot, *userSlot) {
	return &userSlot{
			user:    author,
			status:  slotAuthor,
			manager: manager,
			state: &stream.User{
				IsCamEnabled: true,
				IsMuted:      true,
				UserUUID:     "",
				UserName:     "",
			},
			inputAVConn:    nil,
			outputAVStream: &avStream{stream: nil, closeCh: make(chan struct{})},
			stateStream:    &stateStream{stream: nil, closeCh: make(chan struct{})},
			chatStream:     &chatStream{stream: nil, closeCh: make(chan struct{})},
		},
		&userSlot{
			user:           nil,
			status:         slotFree,
			manager:        manager,
			state:          nil,
			inputAVConn:    nil,
			outputAVStream: &avStream{stream: nil, closeCh: nil},
			stateStream:    &stateStream{stream: nil, closeCh: nil},
			chatStream:     &chatStream{stream: nil, closeCh: nil},
		}
}

func makeChatState() *chatState {
	return &chatState{
		mu: sync.RWMutex{},
		messages: []*chat.ChatMessage{{
			UUID:     uuid.NewString(),
			UserUUID: uuid.NewString(),
			UserName: "Server",
			Time:     0,
			Msg:      "Welcome to chat",
		}},
	}
}
