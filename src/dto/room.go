package dto

import "github.com/google/uuid"

type Room struct {
	UUID   uuid.UUID `json:"uuid"`
	Author *User     `json:"author"`
	Guest  *User     `json:"guest"`
}
