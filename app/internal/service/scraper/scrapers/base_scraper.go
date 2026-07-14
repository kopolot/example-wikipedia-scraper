package scrapers

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/service/browser"
	browserTypes "example-wikipedia-scraper/internal/types/browser"
	types "example-wikipedia-scraper/internal/types/scraper"
	"example-wikipedia-scraper/pkg/helpers"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var expiresTimeSinceEpoch = cdp.TimeSinceEpoch(time.Now().Add(24 * 365 * time.Hour))

type ListPageResponse struct {
	Pages []*ListPageData `json:"pages"`
}

type ListPageData struct {
	URL        string `json:"URL"`
	ExternalID string `json:"externalID"`
	No         int    `json:"no"`
	Page       int    `json:"page"`
}

type BaseScraper struct {
	lastPageId           string
	newLastPageId        string
	listUnformatedBaseUrl string
	browser               interfaces.BrowserInterface
	logger                interfaces.LoggerInterface
	config                *config.SiteConfig
	pagesProcessed       int
	scriptsCache          map[string]string
	browserOptions        *browserTypes.FetchOptions
	scrapeOptions         *types.ScrapeOptions
	listPageDtoChan        chan *ListPageData
	pageDtoChanCloseMutex  sync.Mutex
	fetchMutex            sync.Mutex
	pageDtoChanClosed      bool
	scriptsCacheMutex     sync.Mutex
}

func (s *BaseScraper) ClosePageDtoChan() {
	s.pageDtoChanCloseMutex.Lock()
	defer s.pageDtoChanCloseMutex.Unlock()
	if !s.pageDtoChanClosed {
		close(s.listPageDtoChan)
		s.pageDtoChanClosed = true
	}
}

func (s *BaseScraper) GetConfig() config.SiteConfig {
	return *s.config
}

func (s *BaseScraper) GetName() string {
	return s.config.Name
}

func (s *BaseScraper) GetScriptFileContent(fileName string) (string, error) {
	scrapername := s.GetConfig().Name
	s.scriptsCacheMutex.Lock()
	defer s.scriptsCacheMutex.Unlock()
	if content, found := s.scriptsCache[fileName]; found {
		return content, nil
	}
	filePath := helpers.GetAbsoluteFilePath(filepath.Join("..", "..", "..", "..", "resource", "scraper", "js", scrapername, fileName+".js"))
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", s.logAndReturnError("error reading script file", scrapername, "file", fileName, "err", err, "filePath", filePath)
	}
	s.scriptsCache[fileName] = string(content)
	return string(content), nil
}

func (s *BaseScraper) GetWorkerConfig(name string) config.WorkerConfig {
	for _, worker := range s.GetConfig().Workers {
		if worker.Name == name {
			return *worker
		}
	}
	return config.WorkerConfig{
		Name:     name,
		NumberOf: 1,
		Cooldown: 1000, // Default cooldown in milliseconds
	}
}

func (s *BaseScraper) GetLastPageId(sitename string) (string, error) {
	return s.lastPageId, nil
}

func (s *BaseScraper) getBrowserOptions() browserTypes.FetchOptions {
	s.InitBrowserOptions()
	return *s.browserOptions
}

func (s *BaseScraper) SetBrowserOptions(opts ...browserTypes.FetchOption) {
	s.InitBrowserOptions()
	for _, opt := range opts {
		opt(s.browserOptions)
	}
}

func (s *BaseScraper) SetScrapeOptions(opts ...types.ScrapeOption) {
	s.scrapeOptions = types.ApplyOptions(opts...)
}

func (s *BaseScraper) GetScrapeOptions() *types.ScrapeOptions {
	if s.scrapeOptions == nil {
		s.scrapeOptions = &types.ScrapeOptions{}
	}
	return s.scrapeOptions
}

func (s *BaseScraper) InitBrowserOptions() {
	if s.browserOptions == nil {
		browserOptions := s.browser.GetOptions()
		s.browserOptions = &browserOptions
	}
}

func (s *BaseScraper) InitScraper(opts ...types.ScrapeOption) error {
	s.SetScrapeOptions(opts...)
	// return s.buildBrowserOptions()
	return nil
}

func (s *BaseScraper) buildBrowserOptions() error {
	options := s.GetScrapeOptions()
	timeout := options.Timeout
	if timeout == 0 {
		t := s.getBrowserOptions().Timeout
		timeout = t
	}
	_ = s.GetConfig()
	const ratelimitRetryCooldownMs = 12000
	s.SetBrowserOptions(
		browserTypes.WithTimeout(timeout),
		browserTypes.WithRatelimitCooldown(ratelimitRetryCooldownMs*time.Millisecond),
		browserTypes.WithRetries(7),
		browserTypes.WithSaveBody(true),
		browserTypes.WithCookies(options.Cookies...),
	)
	return nil
}

func (s *BaseScraper) StartListScraperWorker() error {
	siteConfig := s.GetConfig()
	pageCount := siteConfig.PagesBack
loop:
	for page := 1; page <= pageCount; page++ {
		if s.processListScraperPageReturnBreakLoop(page) {
			break loop
		}
	}
	if s.newLastPageId != "" {
		s.lastPageId = s.newLastPageId
	}
	return nil
}

func (s *BaseScraper) processListScraperPageReturnBreakLoop(page int) bool {
	siteConfig := s.GetConfig()
	siteName := siteConfig.Name
	waitBeforeScrapeAttempt(siteName)

	workerConfig := s.GetWorkerConfig("main")
	pageCount := siteConfig.PagesBack
	err := s.ScrapeListPage(page)
	pageUrl := fmt.Sprintf(s.listUnformatedBaseUrl, page)
	switch err {
	case nil:
		recordScrapeSuccess(siteName)
		if page < pageCount {
			s.StartCooldown(workerConfig.Cooldown)
		}
	case types.ErrTargetServer:
		if backoff := recordScrapeFailure(siteName, err); backoff > 0 {
			time.Sleep(backoff)
		}
		if shouldLogScrapeError(siteName) {
			s.log(logger.LevelWarn, "Target server error encountered, retrying later", siteName, "page", pageUrl, "err", err, "options", s.GetScrapeOptions())
		}
	case types.ErrLastPageReached:
		recordScrapeSuccess(siteName)
		s.log(logger.LevelDebug, "Last page reached on page, stop scraping", siteName, "page", pageUrl, "options", s.GetScrapeOptions())
		return true
	default:
		if backoff := recordScrapeFailure(siteName, err); backoff > 0 {
			time.Sleep(backoff)
		}
		if shouldLogScrapeError(siteName) {
			s.log(logger.LevelError, "Error scraping list page", siteName, "page", pageUrl, "err", err, "options", s.GetScrapeOptions())
		}
	}
	return false
}

func (s *BaseScraper) ScrapeListPage(page int) error {
	if s.pageDtoChanClosed {
		return nil
	}
	pageUrl := fmt.Sprintf(s.listUnformatedBaseUrl, page)
	err := s.FetchPageListDataAndSendToDtoToChannel(pageUrl)
	return err
}

func (s *BaseScraper) FetchPageListDataAndSendToDtoToChannel(url string) error {
	listPageResponse, err := s.FetchPageListData(url)
	if err != nil {
		return err
	}
	return s.SendListPagesToChannel(listPageResponse.Pages)
}

func (s *BaseScraper) FetchPageListData(url string) (ListPageResponse, error) {
	var listPageResponse ListPageResponse
	if s.pageDtoChanClosed {
		return listPageResponse, nil
	}
	ctx, cancel, err := s.browser.GetNewContext()
	if err != nil {
		return listPageResponse, s.logAndReturnError("error creating browser context for page list data", s.GetName(), "url", url, "error", err)
	}
	defer cancel()
	browserSession := browser.NewBrowserSession(ctx, cancel)
	browserSession.SetURL(url)
	err = s.fetchPageWithRetry(browserSession)
	if err != nil {
		return listPageResponse, s.translateBrowserError(err)
	}
	return s.GetDataFromList(browserSession.GetContext())
}

func (s *BaseScraper) SendListPagesToChannel(pages []*ListPageData) error {
	siteConfig := s.GetConfig()
	options := s.GetScrapeOptions()
	for _, page := range pages {
		if s.pageDtoChanClosed {
			return nil
		}
		if lastPageId, _ := s.GetLastPageId(siteConfig.Name); lastPageId != "" && page.ExternalID == lastPageId {
			s.log(logger.LevelDebug, "Last page reached, stopping scraping", siteConfig.Name, "pageURL", page.URL)
			return types.ErrLastPageReached
		}
		if options.MaxItems <= 0 || s.pagesProcessed < options.MaxItems {
			s.listPageDtoChan <- page
			if page.Page == 1 && page.No == 1 {
				s.newLastPageId = page.ExternalID
			}
			s.pagesProcessed++
		} else {
			s.ClosePageDtoChan()
			return nil
		}
	}
	return nil
}

func (s *BaseScraper) RunChromedp(ctx context.Context, actions ...chromedp.Action) error {
	return s.browser.RunActions(ctx, actions...)
}

func (s *BaseScraper) GetDataFromList(requestCtx context.Context) (ListPageResponse, error) {
	scriptContent, err := s.GetScriptFileContent("listPage")
	var listPageResponse ListPageResponse
	if err != nil {
		return listPageResponse, s.logAndReturnError("error getting script content", s.GetConfig().Name, "err", err)
	}
	s.RunChromedp(requestCtx,
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second))
	err = s.RunChromedp(requestCtx, chromedp.Evaluate(
		scriptContent,
		&listPageResponse,
	))
	if err != nil {
		return listPageResponse, err
	}
	return listPageResponse, nil
}

func (s *BaseScraper) ScrapeFullDataFromPagesAndSendToChan(channels *types.ScrapeChannels, scrapePageCallback func(string) (*dto.PageDTO, error)) error {
	pageQueue := channels.PageQueue
	workerConfig := s.GetWorkerConfig("details")
	siteConfig := s.GetConfig()
	_ = s.getBrowserOptions()
	_, err := s.GetScriptFileContent("getPageDetails")
	if err != nil {
		s.log(logger.LevelError, "Error getting script content for page details", siteConfig.Name, "err", err)
		return fmt.Errorf("error getting script content for page details: %w", err)
	}
	var wg sync.WaitGroup
	for range workerConfig.NumberOf {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				pageData, ok := <-s.listPageDtoChan
				if !ok {
					break
				}
				s.log(logger.LevelDebug, "Processing page", siteConfig.Name, "pageURL", pageData.URL)
				if pageData.URL == "" {
					s.log(logger.LevelError, "Empty page URL, skipping", siteConfig.Name, "pageURL", pageData.URL)
					continue
				}
				pageDto, err := scrapePageCallback(pageData.URL)
				if err != nil {
					channels.FailedPages <- &dto.UnprocessedPageDTO{
						URL:      pageData.URL,
						SiteName: siteConfig.Name,
					}
					s.log(logger.LevelError, "Error scraping data for page", siteConfig.Name, "pageURL", pageData.URL, "err", err)
					continue
				}
				if pageDto != nil {
					pageQueue <- pageDto
				}
				s.StartCooldown(workerConfig.Cooldown)
			}
		}()
	}
	wg.Wait()
	return nil
}

func (s *BaseScraper) log(level int8, msg string, sitename string, args ...any) {
	newArgs := make([]any, 0, len(args)+2)
	newArgs = append(newArgs, "sitename", sitename)
	newArgs = append(newArgs, args...)
	s.logger.Log(level, msg, newArgs...)
}

func (s *BaseScraper) logAndReturnError(msg string, sitename string, args ...any) error {
	s.log(logger.LevelError, msg, sitename, args...)
	return fmt.Errorf(msg, args...)
}

func (s *BaseScraper) getFormattedErrorMessage(baseMsg string, args ...any) error {
	return fmt.Errorf("%s | args=%v", baseMsg, args)
}

// start cooldown in milliseconds
func (s *BaseScraper) StartCooldown(cooldown int) {
	if cooldown <= 0 {
		return
	}
	jitterSpread := cooldown / 10
	jitter := rand.Intn(jitterSpread) - jitterSpread/2
	time.Sleep(time.Duration(cooldown)*time.Millisecond + time.Duration(jitter)*time.Millisecond)
}

func (s *BaseScraper) ValidatePageBase(page *model.Page, callback func(context.Context) (*dto.PageDTO, error)) (*dto.PageDTO, error) {
	pageCtx, cancel, err := s.browser.GetNewContext()
	if err != nil {
		s.log(logger.LevelError, "Error creating browser context for page validation", s.GetName(), "pageURL", page.URL, "error", err)
		return nil, types.ErrApplication
	}
	defer cancel()
	pageURL := page.URL
	browserSession := browser.NewBrowserSession(pageCtx, cancel)
	browserSession.SetURL(pageURL)
	err = s.fetchPageWithRetry(browserSession)
	if err != nil {
		return nil, s.translateBrowserError(err)
	}
	// tu sie pozbyc context ??
	pageDto, err := callback(browserSession.GetContext())
	return pageDto, s.translateBrowserError(err)
}

func (s *BaseScraper) translateBrowserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, browserTypes.ErrFetchPageNotFound) {
		return types.ErrRecordNotFound
	}
	if errors.Is(err, browserTypes.ErrFetchRatelimit) {
		s.log(logger.LevelDebug, "Translated browser error", s.GetConfig().Name, "originalErr", err, "translatedErr", types.ErrRatelimit)
		return types.ErrRatelimit
	}
	if errors.Is(err, browserTypes.ErrFetchManagedChallenge) {
		s.log(logger.LevelDebug, "Translated browser error", s.GetConfig().Name, "originalErr", err, "translatedErr", types.ErrRatelimit)
		return types.ErrRatelimit
	}
	if errors.Is(err, browserTypes.ErrFetchTargetServer) {
		s.log(logger.LevelDebug, "Translated browser error", s.GetConfig().Name, "originalErr", err, "translatedErr", types.ErrTargetServer)
		return types.ErrTargetServer
	}
	s.log(logger.LevelDebug, "Translated browser error", s.GetConfig().Name, "originalErr", err, "translatedErr", types.ErrApplication)
	return types.ErrApplication
}

func (s *BaseScraper) ScrapeAsyncBase(
	channels *types.ScrapeChannels,
	scrapePageCallback func(string) (*dto.PageDTO, error),
) error {
	workerConfig := s.GetWorkerConfig("main")
	go func() {
		for {
			waitBeforeScrapeAttempt(s.GetName())
			s.StartListScraperWorker()
			s.StartCooldown(workerConfig.Cooldown)
		}
	}()
	return s.ScrapeFullDataFromPagesAndSendToChan(channels, scrapePageCallback)
}

func (s *BaseScraper) ScrapePageDataBase(
	pageURL string,
	getDtoFunc func(context.Context) (*dto.PageDTO, error),
) (*dto.PageDTO, error) {
	ctx, cancel, err := s.browser.GetNewContext()
	if err != nil {
		return nil, s.logAndReturnError("error creating browser context for page validation", s.GetName(), "pageURL", pageURL, "error", err)
	}
	defer cancel()
	browserSession := browser.NewBrowserSession(ctx, cancel)
	browserSession.SetURL(pageURL)
	err = s.fetchPageWithRetry(browserSession)
	if err != nil {
		return nil, s.translateBrowserError(err)
	}
	s.log(logger.LevelDebug, "Fetching page details", s.GetName(), "pageURL", pageURL)
	localPageDto, err := getDtoFunc(browserSession.GetContext())
	if err != nil {
		s.log(logger.LevelError, "Error getting page details from browser context", s.GetName(), "pageURL", pageURL, "error", err)
		return nil, fmt.Errorf("error getting page details from browser context: %w", err)
	}
	return localPageDto, nil
}

func (s *BaseScraper) fetchPageWithRetry(browserSession interfaces.BrowserSessionInterface) error {
	s.fetchMutex.Lock()
	defer s.fetchMutex.Unlock()
	s.buildBrowserOptions()
	browserOptions := s.getBrowserOptions()
	browserSession.SetOptions(browserOptions)
	return s.browser.FetchPageWithRetry(browserSession)
}
