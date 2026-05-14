package auth

import (
	"errors"
	"example-wikipedia-scraper/internal/model"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserNotVerified    = errors.New("user email not verified")
	ErrTokenExpired       = errors.New("token has expired")
)

type AuthManagerInterface interface {
	Authenticate(token string) (*model.User, error)
	GenerateToken(user model.User, expirationTime time.Duration) (string, error)
	Login(login string, password string) (string, error)
}
