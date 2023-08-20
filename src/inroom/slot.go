package inroom

import (
	"errors"
	"golang.org/x/net/websocket"
	"web-grpc-video-chat/src/dto"
	"web-grpc-video-chat/src/inroom/chat"
	"web-grpc-video-chat/src/inroom/stream"
)

type userSlotStatus int

const (
	slotFree     userSlotStatus = 0
	slotOccupied                = 1
	slotAuthor                  = 2
)

type chatStream struct {
	stream  chat.Chat_ListenServer
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
	user           *dto.User
	manager        *RoomManager
	state          *stream.User
	inputAVConn    *websocket.Conn
	outputAVStream *avStream
	stateStream    *stateStream
	chatStream     *chatStream
}

func (us *userSlot) isUser(user *dto.User) bool {
	return us.user != nil && us.user == user
}

func (us *userSlot) takeSlot(user *dto.User) error {
	if us.status != slotFree {
		return errors.New("Slot is occupied. ")
	}
	us.user = user
	us.status = slotOccupied
	us.outputAVStream.closeCh = make(chan struct{})
	us.stateStream.closeCh = make(chan struct{})
	us.chatStream.closeCh = make(chan struct{})
	return nil
}

func (us *userSlot) freeUpSlot() *dto.User {
	var user *dto.User
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

func (us *userSlot) sendChatMessage(message *chat.ChatMessage) {
	if us.status != slotFree {
		us.chatStream.stream.Send(message)
	}
}

func (us *userSlot) connectAVStream(server stream.Stream_AVStreamServer) chan struct{} {
	close(us.outputAVStream.closeCh)
	us.outputAVStream.stream = server
	us.outputAVStream.closeCh = make(chan struct{})
	return us.outputAVStream.closeCh
}

func (us *userSlot) connectChatStream(server chat.Chat_ListenServer) chan struct{} {
	close(us.chatStream.closeCh)
	us.chatStream.stream = server
	us.chatStream.closeCh = make(chan struct{})
	return us.chatStream.closeCh
}

func (us *userSlot) connectStateStream(server stream.Stream_StreamStateServer) chan struct{} {
	close(us.stateStream.closeCh)
	us.stateStream.stream = server
	us.stateStream.closeCh = make(chan struct{})
	return us.stateStream.closeCh
}
