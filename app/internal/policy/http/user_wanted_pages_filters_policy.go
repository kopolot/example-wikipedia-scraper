package http

import (
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/policy"
	"errors"
)

var (
	ErrNoRecordFound    = errors.New("no record found")
	ErrPermissionDenied = errors.New("permission denied")
)

type UserWantedPagesFiltersPolicy struct {
	subscriptionLevelPolicy *policy.SubscriptionLevelPolicy
}

func NewUserWantedPagesFiltersPolicy(
	subscriptionService serviceInterfaces.SubscriptionServiceInterface,
) *UserWantedPagesFiltersPolicy {
	return &UserWantedPagesFiltersPolicy{
		subscriptionLevelPolicy: policy.NewSubscriptionLevelPolicy(subscriptionService),
	}
}

func (p *UserWantedPagesFiltersPolicy) CanUpdate(user *model.User, filter *model.UserWantedPagesFilter) error {
	if user == nil || filter == nil {
		return ErrNoRecordFound
	}
	if user.Role == model.RoleAdmin {
		return nil
	}
	if user.ID != filter.UserID {
		return ErrPermissionDenied
	}
	if err := p.subscriptionLevelPolicy.CanDoPremiumAction(user); err != nil {
		return err
	}
	return nil
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
	if err := p.subscriptionLevelPolicy.CanCreateUserWantedPagesFilter(user); err != nil {
		return err
	}
	return nil
}

func (p *UserWantedPagesFiltersPolicy) CanDelete(user *model.User, filter *model.UserWantedPagesFilter) error {
	return p.CanUpdate(user, filter)
}

func (p *UserWantedPagesFiltersPolicy) CanGetFilteredPages(user *model.User, _ *model.UserWantedPagesFilter) error {
	return p.subscriptionLevelPolicy.CanDoPremiumAction(user)
}
