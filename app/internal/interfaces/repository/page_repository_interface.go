package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
	"time"
)

type PageRepositoryInterface interface {
	repository.RepositoryInterface[*model.Page]
	GetLastFromSite(site string) (*model.Page, error)
	FindByName(name string) ([]*model.Page, error)
	GetPagesUpdatedBefore(t time.Duration) ([]*model.Page, error)
	GetPageAndSimilarPages(page *model.Page) ([]*model.Page, error)
}
