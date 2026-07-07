package http

import (
	repoInterface "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	"errors"
)

const MaxFiltersPerUser = 10

var (
	ErrNoRecordFound    = errors.New("no record found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrFilterLimit      = errors.New("maximum number of filters reached")
)

type UserWantedPagesFiltersPolicy struct {
	filterRepo repoInterface.UserWantedPagesFilterRepositoryInterface
}

func NewUserWantedPagesFiltersPolicy(
	filterRepo repoInterface.UserWantedPagesFilterRepositoryInterface,
) *UserWantedPagesFiltersPolicy {
	return &UserWantedPagesFiltersPolicy{filterRepo: filterRepo}
}

func (p *UserWantedPagesFiltersPolicy) CanView(user *model.User, filter *model.UserWantedPagesFilter) error {
	if user == nil || filter == nil {
		return ErrNoRecordFound
	}
	if user.Role == model.RoleAdmin {
		return nil
	}
	if user.ID != filter.UserID {
		return ErrPermissionDenied
	}
	return nil
}

func (p *UserWantedPagesFiltersPolicy) CanCreate(user *model.User, _ *model.UserWantedPagesFilter) error {
	if user == nil {
		return ErrNoRecordFound
	}
	count, err := p.filterRepo.CountBy("user_id = ?", user.ID)
	if err != nil {
		return err
	}
	if count >= MaxFiltersPerUser {
		return ErrFilterLimit
	}
	return nil
}

func (p *UserWantedPagesFiltersPolicy) CanUpdate(user *model.User, filter *model.UserWantedPagesFilter) error {
	return p.CanView(user, filter)
}

func (p *UserWantedPagesFiltersPolicy) CanDelete(user *model.User, filter *model.UserWantedPagesFilter) error {
	return p.CanView(user, filter)
}

func (p *UserWantedPagesFiltersPolicy) CanGetFilteredPages(user *model.User, _ *model.UserWantedPagesFilter) error {
	if user == nil {
		return ErrNoRecordFound
	}
	return nil
}
