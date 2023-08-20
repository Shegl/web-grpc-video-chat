package services

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"web-grpc-video-chat/src/dto"
)

type AuthService struct {
	authUsers map[uuid.UUID]*dto.User
	repo      *dto.Repository
}

var randomNames = [10]string{"Ron", "John", "Don", "Hubert", "Mike", "Alex", "Anton", "Mathias", "Dora", "Jane"}
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 \x{00F0}-\x{02AF}]+`)

func (a *AuthService) Authenticate(userName string) (*dto.User, error) {
	userName = nonAlphanumericRegex.ReplaceAllString(userName, "")
	if userName == "" {
		userName = randomNames[rand.Intn(len(randomNames))]
	}

	user, err := a.repo.CreateUser(userName)
	if err != nil {
		return a.Authenticate(userName)
	}

	return user, nil
}

func (a *AuthService) GetUser(userUUID uuid.UUID) (*dto.User, error) {
	user := a.repo.FindUserByUuid(userUUID)
	if user == nil {
		return nil, errors.New("There is no such user. ")
	}
	return user, nil
}

// Logout actually we want our users to be logged in forever
// but, it's nice to have option to logout if you have logged in, right ?
func (a *AuthService) Logout(userUUID uuid.UUID) {
	a.repo.ForgetUserByUuid(userUUID)
}

func NewAuthService(repo *dto.Repository) *AuthService {
	return &AuthService{repo: repo}
}
