package testutils

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (s *MockUserService) UpdateUser(user *model.User, dto dto.UpdateUserDTO) (*model.User, error) {
	args := s.Called(user, dto)
	return args.Get(0).(*model.User), args.Error(1)
}

func (s *MockUserService) CreateUser(dto dto.CreateUserDTO) (*model.User, error) {
	args := s.Called(dto)
	return args.Get(0).(*model.User), args.Error(1)
}

func (s *MockUserService) ChangeEmail(user *model.User, dto dto.ChangeEmailDTO) (*model.User, error) {
	args := s.Called(user, dto)
	return args.Get(0).(*model.User), args.Error(1)
}

func (s *MockUserService) ChangePassword(user *model.User, dto dto.ChangePasswordDTO) error {
	args := s.Called(user, dto)
	return args.Error(0)
}

func (s *MockUserService) VerifyEmail(token string) error {
	args := s.Called(token)
	return args.Error(0)
}

func (s *MockUserService) ForgotPassword(dto dto.ForgotPasswordDTO) error {
	args := s.Called(dto)
	return args.Error(0)
}

func (s *MockUserService) ResetPassword(dto dto.ResetPasswordDTO) error {
	args := s.Called(dto)
	return args.Error(0)
}

func (s *MockUserService) Logout(user *model.User) error {
	args := s.Called(user)
	return args.Error(0)
}
