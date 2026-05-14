package interfaces

import (
	"context"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/scraper"
)

type ScraperInterface interface {
	GetName() string
	GetURL() string
	ScrapeSync(opts ...types.ScrapeOption) ([]model.Page, error)
	ScrapeAsync(channels *types.ScrapeChannels, opts ...types.ScrapeOption) error
	ScrapePageData(pageURL string, pageCtx context.Context) (*dto.PageDTO, error)
	ValidatePage(page *model.Page) (*dto.PageDTO, error)
}
