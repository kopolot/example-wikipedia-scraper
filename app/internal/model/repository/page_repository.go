package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
	"time"
)

type PageRepository struct {
	*repository.GenericRepository[*model.Page]
}

func NewPageRepository() *PageRepository {
	return &PageRepository{
		GenericRepository: repository.NewGenericRepository[*model.Page](),
	}
}

func (r *PageRepository) GetLastFromSite(site string) (*model.Page, error) {
	var page model.Page
	err := r.QB.Where("site_name = ?", site).Order("created_at desc").First(&page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) FindByName(name string) ([]*model.Page, error) {
	var models []*model.Page
	err := r.QB.Where("name LIKE ?", "%"+name+"%").Find(&models)
	return models, err
}

func (r *PageRepository) GetPagesUpdatedBefore(t time.Duration) ([]*model.Page, error) {
	var models []*model.Page
	err := r.QB.Where("updated_at <= ?", time.Now().Add(t)).Order("updated_at asc").Find(&models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *PageRepository) GetPageAndSimilarPages(page *model.Page) ([]*model.Page, error) {
	var models []*model.Page
	err := r.QB.Where("hash_key = ? ", page.HashKey).Find(&models)
	if err != nil {
		return nil, err
	}
	return models, nil
}
