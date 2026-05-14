package testutils

import (
	"example-wikipedia-scraper/internal/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockAuthManager struct {
	mock.Mock
}

func (m *MockAuthManager) Authenticate(token string) (*model.User, error) {
	args := m.Called(token)
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockAuthManager) GenerateToken(user model.User, expirationTime time.Duration) (string, error) {
	args := m.Called(user, expirationTime)
	return args.String(0), args.Error(1)
}
func (m *MockAuthManager) Login(login string, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}
