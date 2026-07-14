package service

import (
	repoInterfaces "example-wikipedia-scraper/internal/interfaces/repository"
	interfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	pkgDb "example-wikipedia-scraper/pkg/db"
	pkgRepo "example-wikipedia-scraper/pkg/repository"
	"time"
)

type SubscriptionService struct {
	filterRepo        repoInterfaces.UserWantedPagesFilterRepositoryInterface
	subLvlRepo        pkgRepo.RepositoryInterface[*model.SubscriptionLevel]
	subLvlProductRepo pkgRepo.RepositoryInterface[*model.SubscriptionLevelProduct]
	userRepo          repoInterfaces.UserRepositoryInterface
}

func NewSubscriptionService(
	filterRepo repoInterfaces.UserWantedPagesFilterRepositoryInterface,
	subLvlRepo pkgRepo.RepositoryInterface[*model.SubscriptionLevel],
	subLvlProductRepo pkgRepo.RepositoryInterface[*model.SubscriptionLevelProduct],
	userRepo repoInterfaces.UserRepositoryInterface,
) *SubscriptionService {
	return &SubscriptionService{
		filterRepo:        filterRepo,
		subLvlRepo:        subLvlRepo,
		subLvlProductRepo: subLvlProductRepo,
		userRepo:          userRepo,
	}
}

func (s *SubscriptionService) HasActiveSubscription(user *model.User) (bool, error) {
	if user.SubscriptionLevel == 0 {
		return false, interfaces.ErrDontHaveSubscription
	}
	if user.SubscriptionExpiration != nil && user.SubscriptionExpiration.Before(time.Now()) {
		return false, interfaces.ErrSubscriptionExpired
	}
	return true, nil
}

func (s *SubscriptionService) CanCreateUserWantedPagesFilter(user *model.User) (bool, error) {
	_, err := s.HasActiveSubscription(user)
	if err != nil {
		return false, err
	}
	count, err := s.filterRepo.CountBy("user_id = ?", user.ID)
	if err != nil {
		return false, err
	}
	lvl, err := s.subLvlRepo.FindOneBy("level = ?", user.SubscriptionLevel)
	if err != nil {
		return false, err
	}
	if lvl.Limit > 0 && count >= int64(lvl.Limit) {
		return false, interfaces.ErrLowSubscription
	}
	return true, nil
}

func (s *SubscriptionService) GetRepository() pkgRepo.RepositoryInterface[*model.SubscriptionLevel] {
	return s.subLvlRepo
}

func (s *SubscriptionService) GetSubscriptionLevelProductRepo() pkgRepo.RepositoryInterface[*model.SubscriptionLevelProduct] {
	return s.subLvlProductRepo
}

func (s *SubscriptionService) GetSubscriptionLevelProducts() ([]*model.SubscriptionLevelProduct, error) {
	return s.subLvlProductRepo.GetAllWithPreloads()
}

// napisac testy do tego
func (s *SubscriptionService) AddSubscriptionTime(user *model.User, level int64, timeToAdd time.Duration) error {
	subTime := user.SubscriptionExpiration
	if subTime != nil {
		if subTime.After(time.Now()) && user.SubscriptionLevel != level {
			return interfaces.ErrAddingSubscriptionTimeWithDifferentLevel
		}
		if subTime.Before(time.Now()) {
			timeNow := time.Now()
			subTime = &timeNow
		}
	} else {
		timeNow := time.Now()
		subTime = &timeNow
	}
	subNewTime := subTime.Add(timeToAdd)
	user.SubscriptionExpiration = &subNewTime
	filterCount, err := s.filterRepo.CountBy("user_id = ?", user.ID)
	if err != nil {
		return err
	}
	subLvl, err := s.subLvlRepo.FindOneBy("level = ?", level)
	if err != nil {
		return err
	}
	qb := s.filterRepo.GetQueryBuilder()
	return qb.Transaction(func(tx pkgDb.QueryBuilder) error {
		if subLvl.Limit < int(filterCount) {
			if err := tx.Where("user_id = ?", user.ID).Delete(&model.UserWantedPagesFilter{}); err != nil {
				return err
			}
		}
		user.SubscriptionLevel = level
		return tx.Updates(user)
	})
}

func (s *SubscriptionService) GetUsersWithActiveSubscription(ids []uint) ([]*model.User, error) {
	return s.userRepo.FindBy("id IN ? AND subscription_expiration > NOW() AND subscription_level > 0", ids)
}
