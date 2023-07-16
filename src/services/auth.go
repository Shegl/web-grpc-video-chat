package services

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"sync"
	"web-grpc-video-chat/src/dto"
)

type AuthService struct {
	authUsers map[uuid.UUID]*dto.User
	mu        sync.RWMutex
}

var randomNames = [10]string{"Ron", "John", "Don", "Hubert", "Mike", "Alex", "Anton", "Mathias", "Dora", "Jane"}
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 \x{00F0}-\x{02AF}]+`)

func (a *AuthService) Authenticate(userName string) (*dto.User, error) {
	userName = nonAlphanumericRegex.ReplaceAllString(userName, "")
	if userName == "" {
		userName = randomNames[rand.Intn(len(randomNames))]
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	a.mu.RLock()
	if _, exists := a.authUsers[userUUID]; exists {
		a.mu.RUnlock()
		return a.Authenticate(userName)
	}
	a.mu.RUnlock()

	a.mu.Lock()
	user := &dto.User{
		Name: userName,
		UUID: userUUID,
	}
	a.authUsers[userUUID] = user
	a.mu.Unlock()

	return user, nil
}

func (a *AuthService) GetUser(userUUID uuid.UUID) (*dto.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if user, exists := a.authUsers[userUUID]; !exists {
		return user, errors.New("there is no user with provided UUID")
	} else {
		return user, nil
	}
}

// Logout actually we want our users to be logged in forever
// but, it's nice to have option to logout if you have logged in, right ?
func (a *AuthService) Logout(userUUID uuid.UUID) {
	a.mu.Lock()
	delete(a.authUsers, userUUID)
	a.mu.Unlock()
}

func NewAuthService() *AuthService {
	return &AuthService{authUsers: make(map[uuid.UUID]*dto.User)}
}
