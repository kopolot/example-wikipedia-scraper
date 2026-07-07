package testutils

import (
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

func (m *MockScraper) InitScraper(opts ...scraper.ScrapeOption) error {
	args := m.Called(opts)
	return args.Error(0)
}

func (m *MockScraper) ScrapeSync(opts ...scraper.ScrapeOption) ([]model.Page, error) {
	args := m.Called(opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Page), args.Error(1)
}

func (m *MockScraper) ScrapeAsync(ch *scraper.ScrapeChannels) error {
	args := m.Called(ch)
	return args.Error(0)
}

func (m *MockScraper) ScrapePageData(url string) (*dto.PageDTO, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PageDTO), args.Error(1)
}

func (m *MockScraper) ValidatePage(page *model.Page) (*dto.PageDTO, error) {
	args := m.Called(page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PageDTO), args.Error(1)
}
