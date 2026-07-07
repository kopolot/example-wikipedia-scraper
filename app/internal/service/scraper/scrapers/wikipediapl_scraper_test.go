package scrapers

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/testutils"
	mockService "example-wikipedia-scraper/internal/testutils/service"
	"example-wikipedia-scraper/internal/types/browser"
	types "example-wikipedia-scraper/internal/types/scraper"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/mock"
)

func TestScrapeList_ExploreWillPageOffsetWork(t *testing.T) {
	mockBrowser := mockService.MockBrowser{}
	mockLogger := testutils.MockLogger{}

	siteConfig := &config.SiteConfig{
		Name:      "wikipedia.pl",
		PagesBack: 1,
		URL:       "https://pl.wikipedia.org",
		Enabled:   true,
		Workers: []*config.WorkerConfig{
			{
				Name:     "main",
				NumberOf: 1,
			},
		},
	}

	newCtx, cancelFunc := chromedp.NewContext(context.Background())
	mockBrowser.On("GetOptions", []browser.FetchOption(nil)).Return(browser.FetchOptions{})
	mockBrowser.On("GetNewContext").Return(newCtx, cancelFunc, nil)
	mockBrowser.On("RunActions", mock.Anything, mock.Anything).Return(nil)

	mockBrowserSession := mockService.MockBrowserSession{}
	mockBrowser.On("FetchPageWithRetry", mock.Anything).Return(nil)
	mockBrowserSession.On("GetContext").Return(newCtx)
	mockBrowserSession.On("SetURL", mock.Anything).Return()
	mockBrowserSession.On("SetOptions", mock.Anything).Return()
	mockBrowserSession.On("GetOptions").Return(browser.FetchOptions{SaveBody: true})

	scraper := NewWikipediaPLScraper(siteConfig.URL, &mockBrowser, siteConfig, &mockLogger)
	pageChann := make(chan *dto.PageDTO, 100)
	channels := &types.ScrapeChannels{
		PageQueue:   pageChann,
		FailedPages: make(chan *dto.UnprocessedPageDTO, 100),
	}

	scraper.InitScraper(types.WithMaxItems(1))
	go func() {
		_ = scraper.ScrapeAsync(channels)
	}()

	time.Sleep(2 * time.Second)
	mockBrowser.AssertExpectations(t)
}
