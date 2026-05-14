package service

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
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
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "Verify your email address"
	verificationLink := s.config.GetApiConfig().PublicFrontendHost + "verify-email?token=" + user.EmailVerificationToken
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "Please verify your email by clicking the following link: " + verificationLink
	return s.mailerService.Send(mail)
}

func (s *UserService) sendCreateUserEmail(user *model.User) error {
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "User account created"
	verificationLink := s.config.GetApiConfig().PublicFrontendHost + "verify-email?token=" + user.EmailVerificationToken
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "Please verify your email by clicking the following link: " + verificationLink
	return s.mailerService.Send(mail)
}

func (s *UserService) sendEmailVerifiedEmail(user *model.User) error {
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "Email address verified"
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "Your email address has been successfully verified."
	return s.mailerService.Send(mail)
}

func (s *UserService) sendPasswordResetEmail(user *model.User) error {
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "Password Reset Request"
	resetLink := s.config.GetApiConfig().PublicFrontendHost + "reset_password?token=" + user.PasswordResetToken
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "You can reset your password by clicking the following link: " + resetLink
	return s.mailerService.Send(mail)
}

func (s *UserService) sendPasswordChangedEmail(user *model.User) error {
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "Password Changed"
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "Your password has been successfully changed."
	return s.mailerService.Send(mail)
}

func (s *UserService) sendLogoutEmail(user *model.User) error {
	mail := s.mailerService.NewMail()
	mail.To = []string{user.Email}
	mail.Subject = "Logout Notification"
	// tu można by użyć szablonu czy czegos tego typu
	mail.Body = "You have been successfully logged out."
	return s.mailerService.Send(mail)
}
