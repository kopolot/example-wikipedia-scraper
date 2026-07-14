package service

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/service/mailer"
	"time"

	interfaces "example-wikipedia-scraper/internal/interfaces"
	repoInterface "example-wikipedia-scraper/internal/interfaces/repository"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"

	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/utilities"

	"github.com/google/uuid"
)

type UserService struct {
	repo          repoInterface.UserRepositoryInterface
	mailerService interfaces.MailerInterface
	config        config.ConfigInterface
}

func NewUserService(repo repoInterface.UserRepositoryInterface, mailerService interfaces.MailerInterface, config config.ConfigInterface) *UserService {
	return &UserService{
		repo:          repo,
		mailerService: mailerService,
		config:        config,
	}
}

func (s *UserService) UpdateUser(user *model.User, dto dto.UpdateUserDTO) (*model.User, error) {
	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.Password != "" {
		hashedPassword, err := utilities.HashPassword(dto.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}
	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(dto dto.CreateUserDTO) (*model.User, error) {
	hashedPassword, err := utilities.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:                  dto.Email,
		Password:               hashedPassword,
		Username:               dto.Username,
		Role:                   model.RoleUser,
		EmailVerificationToken: s.GenerateVerificationToken(),
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	err = s.sendCreateUserEmail(user)
	if err != nil {
		_ = s.repo.Delete(user.ID)
		return nil, err
	}
	return user, nil
}

func (s *UserService) VerifyEmail(token string) error {
	model, err := s.repo.GetByEmailVerificationToken(token)
	if err != nil {
		return err
	}
	if model == nil || model.ID == 0 {
		return serviceInterfaces.ErrInvalidVerificationToken
	}
	model.EmailVerified = true
	model.EmailVerificationToken = ""
	err = s.repo.Update(model)
	if err != nil {
		return err
	}
	err = s.sendEmailVerifiedEmail(model)
	return err
}

func (s *UserService) ChangeEmail(user *model.User, dto dto.ChangeEmailDTO) (*model.User, error) {
	if !s.checkPasswordHash(dto.Password, user.Password) {
		return nil, serviceInterfaces.ErrInvalidPassword
	}
	user.Email = dto.NewEmail
	user.EmailVerified = false
	verificationToken := s.GenerateVerificationToken()
	user.EmailVerificationToken = verificationToken
	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	err = s.Logout(user)
	if err != nil {
		return nil, err
	}
	err = s.sendVerificationEmail(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ChangePassword(user *model.User, dto dto.ChangePasswordDTO) error {
	if !s.checkPasswordHash(dto.OldPassword, user.Password) {
		return serviceInterfaces.ErrInvalidPassword
	}
	hashedPassword, err := utilities.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = s.repo.Update(user)
	if err != nil {
		return err
	}
	err = s.Logout(user)
	if err != nil {
		return err
	}
	err = s.sendPasswordChangedEmail(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ForgotPassword(dto dto.ForgotPasswordDTO) error {
	user, err := s.repo.GetByEmail(dto.Email)
	if err != nil {
		return err
	}
	if user == nil || user.ID == 0 {
		return nil
	}
	token := s.GenerateVerificationToken()
	user.PasswordResetToken = token
	atTime := time.Now().Add(1 * time.Hour)
	user.PasswordResetTokenExpiresAt = &atTime
	err = s.repo.Update(user)
	if err != nil {
		return err
	}
	err = s.sendPasswordResetEmail(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(dto dto.ResetPasswordDTO) error {
	user, err := s.repo.GetByPasswordResetToken(dto.Token)
	if err != nil {
		return err
	}
	if user == nil || user.ID == 0 {
		return serviceInterfaces.ErrInvalidResetToken
	}
	if user.PasswordResetTokenExpiresAt.Before(time.Now()) {
		return serviceInterfaces.ErrExpiredResetToken
	}
	hashedPassword, err := utilities.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.PasswordResetTokenExpiresAt = nil
	err = s.repo.Update(user)
	if err != nil {
		return err
	}
	err = s.Logout(user)
	if err != nil {
		return err
	}
	err = s.sendPasswordChangedEmail(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GenerateVerificationToken() string {
	return uuid.NewString()
}

func (s *UserService) Logout(user *model.User) error {
	err := s.repo.UpdateLastLogout(user)
	return err
}

func (s *UserService) hashPassword(password string) (string, error) {
	return utilities.HashPassword(password)
}

func (s *UserService) checkPasswordHash(password, hash string) bool {
	return utilities.CheckPasswordHash(password, hash)
}

func (s *UserService) sendVerificationEmail(user *model.User) error {
	verificationLink := s.config.GetApiConfig().PublicFrontendHost + "verify-email?token=" + user.EmailVerificationToken
	mail, err := mailer.NewTemplateBuilder().VerificationEmail(user.Email, verificationLink)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}

func (s *UserService) sendCreateUserEmail(user *model.User) error {
	verificationLink := s.config.GetApiConfig().PublicFrontendHost + "verify-email?token=" + user.EmailVerificationToken
	mail, err := mailer.NewTemplateBuilder().CreateUserEmail(user.Email, verificationLink)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}

func (s *UserService) sendEmailVerifiedEmail(user *model.User) error {
	mail, err := mailer.NewTemplateBuilder().EmailVerifiedEmail(user.Email)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}

func (s *UserService) sendPasswordResetEmail(user *model.User) error {
	resetLink := s.config.GetApiConfig().PublicFrontendHost + "reset_password?token=" + user.PasswordResetToken
	mail, err := mailer.NewTemplateBuilder().PasswordResetEmail(user.Email, resetLink)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}

func (s *UserService) sendPasswordChangedEmail(user *model.User) error {
	mail, err := mailer.NewTemplateBuilder().PasswordChangedEmail(user.Email)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}

func (s *UserService) sendLogoutEmail(user *model.User) error {
	mail, err := mailer.NewTemplateBuilder().LogoutEmail(user.Email)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}
