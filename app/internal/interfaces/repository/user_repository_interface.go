package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
)

type UserRepositoryInterface interface {
	repository.RepositoryInterface[*model.User]
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmailVerificationToken(token string) (*model.User, error)
	UpdateLastLogin(user *model.User) error
	FindByLogin(login string) (*model.User, error)
	GetByPasswordResetToken(token string) (*model.User, error)
	UpdateLastLogout(user *model.User) error
}
