package interfaces

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/scraper"
)

type ScraperInterface interface {
	GetName() string
	GetURL() string
	InitScraper(opts ...types.ScrapeOption) error
	ScrapeSync(opts ...types.ScrapeOption) ([]model.Page, error)
	ScrapeAsync(channels *types.ScrapeChannels) error
	ScrapePageData(pageURL string) (*dto.PageDTO, error)
	ValidatePage(page *model.Page) (*dto.PageDTO, error)
}
