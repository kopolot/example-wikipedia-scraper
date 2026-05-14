package service

import (
	"errors"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
)

var (
	ErrInvalidPassword          = errors.New("invalid password")
	ErrInvalidVerificationToken = errors.New("invalid verification token")
	ErrInvalidResetToken        = errors.New("invalid reset token")
	ErrExpiredResetToken        = errors.New("expired reset token")
)

type UserServiceInterface interface {
	UpdateUser(user *model.User, dto dto.UpdateUserDTO) (*model.User, error)
	CreateUser(dto dto.CreateUserDTO) (*model.User, error)
	ChangeEmail(user *model.User, dto dto.ChangeEmailDTO) (*model.User, error)
	ChangePassword(user *model.User, dto dto.ChangePasswordDTO) error
	ForgotPassword(dto dto.ForgotPasswordDTO) error
	VerifyEmail(token string) error
	ResetPassword(dto dto.ResetPasswordDTO) error
	Logout(user *model.User) error
}
