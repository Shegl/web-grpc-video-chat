package dto

import "github.com/google/uuid"

type User struct {
	Name string    `json:"username"`
	UUID uuid.UUID `json:"uuid"`
}
