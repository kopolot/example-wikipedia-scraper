package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
	"time"
)

type UserRepository struct {
	*repository.GenericRepository[*model.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		GenericRepository: repository.NewGenericRepository[*model.User](),
	}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	result, err := r.FindOneBy("email = ?", email)
	return result, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	result, err := r.FindOneBy("username = ?", username)
	return result, err
}

func (r *UserRepository) GetByEmailVerificationToken(token string) (*model.User, error) {
	result, err := r.FindOneBy("email_verification_token = ?", token)
	return result, err
}

func (r *UserRepository) UpdateLastLogin(user *model.User) error {
	user.LastLoginAt = time.Now()
	err := r.Update(user)
	return err
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	result, err := r.FindOneBy("username = ? OR email = ?", login, login)
	return result, err
}

func (r *UserRepository) GetByPasswordResetToken(token string) (*model.User, error) {
	result, err := r.FindOneBy("password_reset_token = ?", token)
	return result, err
}

func (r *UserRepository) UpdateLastLogout(user *model.User) error {
	now := time.Now()
	user.LastLogoutAt = &now
	return r.Update(user)
}
