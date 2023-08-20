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
	authorSlot, guestSlot := makeSlots(room.Author)
	manager := &RoomManager{
		room:    room,
		isAlive: true,
		mu:      sync.RWMutex{},
		roomCtx: roomCtx,
		close:   cancelFunc,
		author:  authorSlot,
		guest:   guestSlot,
		chat:    makeChatState(),
	}
	return manager
}

func makeSlots(author *dto.User) (*userSlot, *userSlot) {
	return &userSlot{
			user:   author,
			status: slotAuthor,
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
