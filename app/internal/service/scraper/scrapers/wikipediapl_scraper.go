package scrapers

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/scraper"
	"fmt"

	"github.com/chromedp/chromedp"
)

const (
	WikipediaPLPagesUri     = "/w/index.php?title=Specjalna:Szukaj&limit=100&offset=%d&ns0=1&sort=create_timestamp_desc&search=w+OR+z+OR+u+OR+o+OR+i+OR+a"
	wikipediaPlLimitPerPage = 100
)

type WikipediaPLScraper struct {
	BaseScraper
	url string
}

func NewWikipediaPLScraper(url string, browser interfaces.BrowserInterface, config *config.SiteConfig, loggerInstance interfaces.LoggerInterface) *WikipediaPLScraper {
	// tu page sprawoicz
	baseListUrl := url + WikipediaPLPagesUri
	s := &WikipediaPLScraper{
		url: url,
		BaseScraper: BaseScraper{
			scriptsCache:          make(map[string]string),
			listPageDtoChan:       make(chan *ListPageData, 1000),
			pageDtoChanClosed:     false,
			browser:               browser,
			pagesProcessed:        0,
			config:                config,
			listUnformatedBaseUrl: baseListUrl,
			logger:                loggerInstance,
		},
	}

	s.log(logger.LevelInfo, "Initialized WikipediaPLScraper", s.GetName())
	return s
}

func (s *WikipediaPLScraper) GetURL() string {
	return s.url
}

func (s *WikipediaPLScraper) ScrapeSync(opts ...types.ScrapeOption) ([]model.Page, error) {
	return nil, nil
}

func (s *WikipediaPLScraper) GetOfferDtoFromBrowserContext(pageCtx context.Context) (*dto.PageDTO, error) {
	var localPageData *dto.PageDTO
	scriptContent, err := s.GetScriptFileContent("getPageDetails")
	if err != nil {
		return nil, fmt.Errorf("error reading script file: %w", err)
	}
	err = chromedp.Run(pageCtx,
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(
			scriptContent,
			&localPageData,
		))
	if err != nil {
		return nil, fmt.Errorf("error evaluating script for page details: %w", err)
	}
	return localPageData, nil
}

// To jest wbrew DRY ale na chama przekopiowałem kod z innego projektu
func (s *WikipediaPLScraper) ScrapeAsync(channels *types.ScrapeChannels, opts ...types.ScrapeOption) error {
	s.InitListScraperWorker(opts...)
	workerConfig := s.GetWorkerConfig("main")
	go func() {
		for {
			s.StartListScraperWorker(opts...)
			s.StartCooldown(workerConfig.Cooldown)
		}
	}()
	return s.ScrapeFullDataFromPagesAndSendToChan(channels, s.ScrapePageData)
}

func (s *WikipediaPLScraper) StartListScraperWorker(opts ...types.ScrapeOption) error {
	siteConfig := s.GetConfig()
	pageCount := siteConfig.PagesBack
	browserCtx, cancel, err := s.browser.GetNewContext()
	if err != nil {
		return s.logAndReturnError("error creating new browser context", siteConfig.Name, "err", err)
	}
	defer cancel()
loop:
	for page := 1; page <= pageCount; page++ {
		offset := (page - 1) * wikipediaPlLimitPerPage
		if s.processListScraperPageReturnBreakLoop(offset, browserCtx, opts...) {
			break loop
		}
	}
	if s.newLastPageId != "" {
		s.lastPageId = s.newLastPageId
	}
	return nil
}

func (s *WikipediaPLScraper) ScrapePageData(pageURL string, pageCtx context.Context) (*dto.PageDTO, error) {
	return s.ScrapePageDataBase(pageURL, pageCtx, s.GetOfferDtoFromBrowserContext)
}

func (s *WikipediaPLScraper) ValidatePage(page *model.Page) (*dto.PageDTO, error) {
	return s.ValidatePageBase(page, s.GetOfferDtoFromBrowserContext)
}

func (s *WikipediaPLScraper) ScrapeListPage(page int, browserCtx context.Context, opts ...types.ScrapeOption) error {
	options := types.ApplyOptions(opts...)
	if s.pageDtoChanClosed {
		return nil
	}
	pageUrl := fmt.Sprintf(s.listUnformatedBaseUrl, (page-1)*100)
	err := s.FetchPageListDataAndSendToDtoToChannel(pageUrl, browserCtx, *options)
	return err
}
