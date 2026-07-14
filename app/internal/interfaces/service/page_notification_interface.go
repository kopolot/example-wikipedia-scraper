package service

import "example-wikipedia-scraper/internal/model"

type PageNotificationServiceInterface interface {
	EnqueuePagesNotifications() error
	RegisterQueueHandlers() error
	GetNotNotifiedPages() ([]*model.Page, error)
	MarkPageAndSimilarAsNotified(page *model.Page) error
	GetMatchedPageFilters(page model.Page) ([]*model.UserWantedPagesFilter, error)
	GetUsersToNotifyByFilters(filters []*model.UserWantedPagesFilter) ([]*model.User, error)
	SendNotificationToUsersAboutPage(emails []string, page model.Page) error
	NotifyUsersAboutPage(usersEmail []string, pageID uint) error
}
