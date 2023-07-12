package dto

import "github.com/google/uuid"

type Room struct {
	State  int       `json:"state"`
	UUID   uuid.UUID `json:"uuid"`
	Author *User     `json:"author"`
	Guest  *User     `json:"guest"`
}
