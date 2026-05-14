package testutils

import (
	"context"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/types/scraper"

	"github.com/stretchr/testify/mock"
)

type MockScraper struct {
	mock.Mock
}

func (m *MockScraper) GetName() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockScraper) GetURL() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockScraper) ScrapeSync(opts ...scraper.ScrapeOption) ([]model.Page, error) {
	args := m.Called(opts)
	return args.Get(0).([]model.Page), args.Error(1)
}
func (m *MockScraper) ScrapeAsync(ch *scraper.ScrapeChannels, opts ...scraper.ScrapeOption) error {
	args := m.Called(ch, opts)
	return args.Error(0)
}
func (m *MockScraper) ScrapePageData(url string, ctx context.Context) (*dto.PageDTO, error) {
	args := m.Called(url, ctx)
	return args.Get(0).(*dto.PageDTO), args.Error(1)
}
func (m *MockScraper) ValidatePage(page *model.Page) (*dto.PageDTO, error) {
	args := m.Called(page)
	return args.Get(0).(*dto.PageDTO), args.Error(1)
}
