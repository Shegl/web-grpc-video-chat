package domain

import "github.com/google/uuid"

type ChatMessage struct {
	Time     int64
	UUID     uuid.UUID
	UserUUID uuid.UUID
	UserName string
	Msg      string
}

func (msg ChatMessage) GetTime() int64 {
	return msg.Time
}

func (msg ChatMessage) GetUUID() uuid.UUID {
	return msg.UUID
}

func (msg ChatMessage) GetUserUUID() uuid.UUID {
	return msg.UserUUID
}

func (msg ChatMessage) GetUserName() string {
	return msg.UserName
}

func (msg ChatMessage) GetMsg() string {
	return msg.Msg
}
