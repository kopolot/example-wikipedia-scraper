package service

import (
	"example-wikipedia-scraper/internal/model"
	pkgRepo "example-wikipedia-scraper/pkg/repository"
	"errors"
	"time"
)

var (
	ErrDontHaveSubscription                     = errors.New("user does not have an active subscription")
	ErrLowSubscription                          = errors.New("user subscription level does not allow this action")
	ErrSubscriptionExpired                      = errors.New("user subscription has expired")
	ErrAddingSubscriptionTimeWithDifferentLevel = errors.New("cannot add subscription time with different level than the one user currently has")
)

type SubscriptionServiceInterface interface {
	HasActiveSubscription(user *model.User) (bool, error)
	CanCreateUserWantedPagesFilter(user *model.User) (bool, error)
	GetRepository() pkgRepo.RepositoryInterface[*model.SubscriptionLevel]
	GetSubscriptionLevelProducts() ([]*model.SubscriptionLevelProduct, error)
	AddSubscriptionTime(user *model.User, level int64, timeToAdd time.Duration) error
	GetSubscriptionLevelProductRepo() pkgRepo.RepositoryInterface[*model.SubscriptionLevelProduct]
	GetUsersWithActiveSubscription(userIds []uint) ([]*model.User, error)
}
