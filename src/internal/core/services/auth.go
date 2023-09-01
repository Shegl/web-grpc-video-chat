package services

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"web-grpc-video-chat/src/internal/core/domain"
	"web-grpc-video-chat/src/internal/core/repo"
)

type AuthService struct {
	authUsers map[uuid.UUID]*domain.User
	repo      *repo.Repository
}

var randomNames = [10]string{"Ron", "John", "Don", "Hubert", "Mike", "Alex", "Anton", "Mathias", "Dora", "Jane"}
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 \x{00F0}-\x{02AF}]+`)

func (a *AuthService) Authenticate(userName string) (*domain.User, error) {
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

func (a *AuthService) GetUser(userUuid uuid.UUID) (*domain.User, error) {
	user := a.repo.FindUserByUuid(userUuid)
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

func NewAuthService(repo *repo.Repository) *AuthService {
	return &AuthService{repo: repo}
}
