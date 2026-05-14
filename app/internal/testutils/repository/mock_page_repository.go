package testutils

import (
	"example-wikipedia-scraper/internal/model"
	"time"
)

type MockPageRepository struct {
	MockGenericRepository[*model.Page]
}

func (m *MockPageRepository) FindByName(name string) ([]*model.Page, error) {
	args := m.Called(name)
	return args.Get(0).([]*model.Page), args.Error(1)
}

func (m *MockPageRepository) GetLastFromSite(site string) (*model.Page, error) {
	args := m.Called(site)
	return args.Get(0).(*model.Page), args.Error(1)
}

func (m *MockPageRepository) GetPagesUpdatedBefore(t time.Duration) ([]*model.Page, error) {
	args := m.Called(t)
	return args.Get(0).([]*model.Page), args.Error(1)
}

func (m *MockPageRepository) GetPageAndSimilarPages(page *model.Page) ([]*model.Page, error) {
	args := m.Called(page)
	return args.Get(0).([]*model.Page), args.Error(1)
}
