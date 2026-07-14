package policy

import (
	interfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
)

type SubscriptionLevelPolicy struct {
	subscriptionService interfaces.SubscriptionServiceInterface
}

func NewSubscriptionLevelPolicy(subscriptionService interfaces.SubscriptionServiceInterface) *SubscriptionLevelPolicy {
	return &SubscriptionLevelPolicy{
		subscriptionService: subscriptionService,
	}
}

func (p *SubscriptionLevelPolicy) CanAccessPremiumContent(user *model.User) error {
	return p.CanDoPremiumAction(user)
}

func (p *SubscriptionLevelPolicy) CanDoPremiumAction(user *model.User) error {
	_, err := p.subscriptionService.HasActiveSubscription(user)
	return err
}

func (p *SubscriptionLevelPolicy) CanCreateUserWantedPagesFilter(user *model.User) error {
	_, err := p.subscriptionService.CanCreateUserWantedPagesFilter(user)
	return err
}
