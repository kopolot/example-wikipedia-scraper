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

var (
	wikimockBrowser mockService.MockBrowser
	wikimockConfig  testutils.MockConfig
	wikilogger      testutils.MockLogger
)

func setUpWikiMocks(t *testing.T) {
	t.Helper()
	wikimockBrowser = mockService.MockBrowser{}
	wikimockConfig = testutils.MockConfig{}
	wikilogger = testutils.MockLogger{}
}

func TestScrapeList_ExploreWillPageOffsetWork(t *testing.T) {

	setUpWikiMocks(t)
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

	testedUrl := "https://pl.wikipedia.org/w/index.php?title=Specjalna:Szukaj&limit=100&offset=0&ns0=1&sort=create_timestamp_desc&search=w+OR+z+OR+u+OR+o+OR+i+OR+a"

	wikimockBrowser.On("GetOptions", []browser.FetchOption(nil)).Return(browser.FetchOptions{})

	newCtx, cancleFunc := chromedp.NewContext(context.Background())
	wikimockBrowser.On("GetNewContext").Return(newCtx, cancleFunc, nil)

	mockBrowserSession := testutils.MockBrowserSession{}

	wikimockBrowser.On("FetchPageWithRetry", newCtx, testedUrl, mock.AnythingOfType("browser.FetchOptions")).Return(&mockBrowserSession, nil)

	mockBrowserSession.On("GetContext").Return(newCtx)

	scraper := NewWikipediaPLScraper(siteConfig.URL, &wikimockBrowser, siteConfig, &wikilogger)
	pageChann := make(chan *dto.PageDTO, 100)
	channels := &types.ScrapeChannels{
		PageQueue:   pageChann,
		FailedPages: make(chan *dto.UnprocessedPageDTO, 100),
	}

	scraper.ScrapeAsync(channels, types.WithMaxItems(1))

	time.Sleep(10 * time.Second) // dajemy czas na wykonanie goroutine

	wikimockBrowser.AssertExpectations(t)
}
