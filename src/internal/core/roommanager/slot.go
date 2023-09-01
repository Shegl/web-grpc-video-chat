package roommanager

import (
	"errors"
	"golang.org/x/net/websocket"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/pb/stream"
)

type userSlotStatus int

const (
	slotFree     userSlotStatus = 0
	slotOccupied                = 1
	slotAuthor                  = 2
)

type chatStream struct {
	stream  ChatListenerServer
	closeCh chan struct{}
}

type stateStream struct {
	stream  stream.Stream_StreamStateServer
	closeCh chan struct{}
}

type avStream struct {
	stream  stream.Stream_AVStreamServer
	closeCh chan struct{}
}

type userSlot struct {
	status         userSlotStatus
	user           *domain.User
	state          *stream.User
	inputAVConn    *websocket.Conn
	outputAVStream *avStream
	stateStream    *stateStream
	chatStream     *chatStream
}

func (us *userSlot) isUser(user *domain.User) bool {
	return us.user != nil && us.user == user
}

func (us *userSlot) takeSlot(user *domain.User) {
	if us.status != slotFree {
		panic(errors.New("Slot is occupied. Abnormal behavior consistency/atomicity violation. "))
	}
	us.user = user
	us.status = slotOccupied
	us.outputAVStream.closeCh = make(chan struct{})
	us.stateStream.closeCh = make(chan struct{})
	us.chatStream.closeCh = make(chan struct{})
}

func (us *userSlot) freeUpSlot() *domain.User {
	var user *domain.User
	if us.status != slotFree {
		user = us.user
		us.user = nil
		us.status = slotFree
		close(us.outputAVStream.closeCh)
		close(us.stateStream.closeCh)
		close(us.chatStream.closeCh)
	}
	return user
}

func (us *userSlot) updateUserState(state *stream.User) {
	us.state = state
}

func (us *userSlot) sendChatMessage(message domain.ChatMessage) {
	if us.status != slotFree {
		us.chatStream.stream.Send(message)
	}
}

func (us *userSlot) setAVStream(server stream.Stream_AVStreamServer) chan struct{} {
	close(us.outputAVStream.closeCh)
	us.outputAVStream.stream = server
	us.outputAVStream.closeCh = make(chan struct{})
	return us.outputAVStream.closeCh
}

func (us *userSlot) setChatStream(server ChatListenerServer) chan struct{} {
	close(us.chatStream.closeCh)
	us.chatStream.stream = server
	us.chatStream.closeCh = make(chan struct{})
	return us.chatStream.closeCh
}

func (us *userSlot) setStateStream(server stream.Stream_StreamStateServer) chan struct{} {
	close(us.stateStream.closeCh)
	us.stateStream.stream = server
	us.stateStream.closeCh = make(chan struct{})
	return us.stateStream.closeCh
}

func makeSlot(user *domain.User) *userSlot {
	var slot *userSlot
	if user != nil {
		slot = &userSlot{
			user:   user,
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
		}
	} else {
		slot = &userSlot{
			user:           nil,
			status:         slotFree,
			state:          nil,
			inputAVConn:    nil,
			outputAVStream: &avStream{stream: nil, closeCh: nil},
			stateStream:    &stateStream{stream: nil, closeCh: nil},
			chatStream:     &chatStream{stream: nil, closeCh: nil},
		}
	}
	return slot
}
