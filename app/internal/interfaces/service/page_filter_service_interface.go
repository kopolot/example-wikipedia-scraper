package service

import "example-wikipedia-scraper/internal/model"

type PageFilterServiceInterface interface {
	GetPagesByFilters(ids []uint, limit, page uint) ([]*model.Page, error)
	FindMatchingFiltersForPage(page model.Page) ([]*model.UserWantedPagesFilter, error)
}
