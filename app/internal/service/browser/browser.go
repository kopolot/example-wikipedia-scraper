package browser

// refractor i solid principles

import (
	"context"
	"errors"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	types "example-wikipedia-scraper/internal/types/browser"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ContextKey string

const (
	AllowRequestResources ContextKey = "allowRequestResources"
)

type Browser struct {
	config       *config.BrowserSettings
	allocCtx     context.Context
	allocCancel  context.CancelFunc
	browserMutex sync.Mutex
	initialized  bool
	logger       interfaces.LoggerInterface
}

func NewBrowser(cfg *config.BrowserSettings, log interfaces.LoggerInterface) *Browser {
	return &Browser{
		config: cfg,
		logger: log,
	}
}

func (b *Browser) GetConfig() config.BrowserSettings {
	return *b.config
}

func (b *Browser) SetConfig(config config.BrowserSettings) {
	b.config = &config
}

func (b *Browser) InitBrowser() error {
	if b.initialized {
		return nil
	}
	b.browserMutex.Lock()
	defer b.browserMutex.Unlock()
	ctx := context.Background()
	opts := make([]chromedp.ExecAllocatorOption, 0)
	opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)
	for k, v := range b.config.GetBrowserEngineSettings() {
		opts = append(opts, chromedp.Flag(k, v))
	}
	b.allocCtx, b.allocCancel = chromedp.NewExecAllocator(ctx, opts...)
	b.initialized = true
	return nil
}

func (b *Browser) getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func (b *Browser) GetNewContext() (context.Context, context.CancelFunc, error) {
	if !b.initialized {
		return nil, nil, types.ErrInitBrowser
	}
	ctx := b.allocCtx
	ctx, cancel := chromedp.NewContext(ctx)
	return ctx, cancel, nil
}

func (b *Browser) GetPageContent(url string, opts ...types.FetchOption) (*types.BrowserResponse, error) {
	if !b.initialized {
		return nil, types.ErrInitBrowser
	}
	taskCtx, cancel, err := b.GetNewContext()
	if err != nil {
		b.logger.Error("could not create new browser context", "err", err)
		return nil, err
	}
	defer cancel()
	options := b.GetOptions(opts...)
	options.SaveBody = true
	browserSession := b.getBrowserSession(taskCtx, options)
	defer browserSession.GetCancelFunc()()
	response, err := b.runChromeDpWithActions(browserSession, url, options)
	if err != nil {
		b.logger.Error("error fetching page", "url", url, "err", err)
		return nil, err
	}
	b.logger.Debug("Fetched page content", "url", url, "status", response.GetStatusCode(), "responseTime", response.GetResponseTime())
	return response, nil
}

func (b *Browser) FetchPage(ctx context.Context, url string, options types.FetchOptions) (interfaces.BrowserSessionInterface, error) {
	if !b.initialized {
		b.logger.Error("browser not initialized")
		return nil, types.ErrInitBrowser
	}
	browserSession := b.getBrowserSession(ctx, options)
	_, err := b.runChromeDpWithActions(browserSession, url, options)
	if err != nil {
		b.logger.Error("error fetching page", "url", url, "err", err)
		return nil, err
	}
	b.logger.Debug("Fetched page content", "url", url)
	return browserSession, nil
}

func (b *Browser) FetchPageWithRetry(parentCtx context.Context, url string, options types.FetchOptions) (interfaces.BrowserSessionInterface, error) {
	if !b.initialized {
		return nil, types.ErrInitBrowser
	}
	var lastErr error
	for range options.Retries {
		browserSession := b.getBrowserSession(parentCtx, options)
		response, err := b.runChromeDpWithActions(browserSession, url, options)
		err = b.handleFetchResult(response, url, err)
		switch {
		case errors.Is(err, types.ErrFetchPage):
			time.Sleep(1 * time.Second)
		case errors.Is(err, types.ErrFetchRatelimit):
			time.Sleep(options.RatelimitCooldown)
		case errors.Is(err, types.ErrFetchTargetServer):
			time.Sleep(1 * time.Second)
		case errors.Is(err, nil):
			return browserSession, nil
		case errors.Is(err, types.ErrFetchPageNotFound):
			return nil, types.ErrFetchPageNotFound
		default:
			b.logger.Warn("unexpected error fetching page, retrying...", "url", url, "err", err)
			time.Sleep(1 * time.Second)
			err = types.ErrFetchPage
		}
		lastErr = err
	}
	return nil, lastErr
}

func (b *Browser) handleFetchResult(response *types.BrowserResponse, url string, err error) error {
	newErr := b.handleFetchError(response, url, err)
	if newErr != nil {
		return newErr
	}
	newErr = b.handleFetchStatusCode(response)
	if newErr != nil {
		return newErr
	}
	return nil
}

func (b *Browser) handleFetchError(response *types.BrowserResponse, url string, err error) error {
	switch {
	case response == nil:
		b.logger.Warn("no response received", "url", url, "err", err)
		return err
	case errors.Is(err, types.ErrFetchPage):
		b.logger.Warn("error fetching page, retrying...", "response_code", response.GetStatusCode(), "response_url", response.GetURL(), "err", err)
		return err
	case err == nil:
		b.logger.Debug("Fetched page content", "response_code", response.GetStatusCode(), "response_url", response.GetURL())
	default:
		b.logger.Warn("unexpected error fetching page", "response_code", response.GetStatusCode(), "response_url", response.GetURL(), "err", err)
		return err
	}
	return nil
}

func (b *Browser) handleFetchStatusCode(response *types.BrowserResponse) error {
	code := response.GetStatusCode()
	switch {
	case code >= 200 && code < 300:
		return nil
	case code == 404 || code == 410:
		return types.ErrFetchPageNotFound
	case code == 403 || code == 429:
		b.logger.Warn("ratelimit reached, waiting before retry", "url", response.URL, "status", response.GetStatusCode())
		return types.ErrFetchRatelimit
	case code >= 500 && code < 600:
		b.logger.Warn("unexpected status code, retrying...", "url", response.URL, "status", response.GetStatusCode())
		return types.ErrFetchTargetServer
	default:
		b.logger.Warn("client error, not retrying", "url", response.URL, "status", response.GetStatusCode())
		return types.ErrFetchPage
	}
}

func (b *Browser) getBrowserSession(ctx context.Context, options types.FetchOptions) interfaces.BrowserSessionInterface {
	ctx, cancel := b.getRequestContext(ctx, options)
	browserSession := NewBrowserSession(ctx, cancel)
	return browserSession
}

func (b *Browser) getRequestContext(ctx context.Context, options types.FetchOptions) (context.Context, context.CancelFunc) {
	var firsrCancel context.CancelFunc
	if options.Timeout > 0 {
		ctx, firsrCancel = context.WithTimeout(ctx, options.Timeout)
	} else {
		timeout := time.Duration(b.GetConfig().Timeout) * time.Millisecond
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ctx, firsrCancel = context.WithTimeout(ctx, timeout)
	}
	ctx, secondCancel := chromedp.NewContext(ctx)
	cancel := context.CancelFunc(func() {
		firsrCancel()
		secondCancel()
	})
	return ctx, cancel
}

func (b *Browser) GetOptions(opts ...types.FetchOption) types.FetchOptions {
	timeout := b.GetConfig().Timeout
	options := &types.FetchOptions{
		Timeout:     time.Duration(timeout) * time.Millisecond,
		SaveBody:    false,
		SaveHeaders: true,
		UserAgent:   "",
		AllowAssets: false,
		Retries:     3,
	}
	for _, opt := range opts {
		opt(options)
	}
	return *options
}

func (b *Browser) runChromeDpWithActions(browserSession interfaces.BrowserSessionInterface, url string, options types.FetchOptions) (*types.BrowserResponse, error) {
	userAgent := b.resolveUserAgent(options)
	actionList := b.buildActionList(url, userAgent, options)
	ctx := browserSession.GetContext()
	session, ok := browserSession.(*BrowserSession)
	if !ok {
		return nil, errors.New("invalid browser session type")
	}
	session.SetURL(url)
	session.SetUserAgent(userAgent)
	session.SetOptions(options.SaveBody, options.SaveHeaders, options.AllowAssets)
	ctx = b.prepareContext(ctx, options)
	var wg sync.WaitGroup
	wg.Add(1)
	session.SetupResponseListener(&wg)
	response, err := chromedp.RunResponse(ctx, actionList...)
	browserSession.SetNetworkResponseAndClearBrowserResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error running chromedp actions: %w, url=%s", err, url)
	}
	wg.Wait()
	browserResponse := browserSession.GetBrowserResponse()
	return &browserResponse, nil
}

func (b *Browser) buildActionList(url, userAgent string, options types.FetchOptions) []chromedp.Action {
	actions := []chromedp.Action{
		network.Enable(),
	}
	if !options.AllowAssets {
		actions = append(actions, fetch.Enable())
	}
	if len(options.Cookies) > 0 {
		actions = append(actions, network.SetCookies(options.Cookies))
	}
	actions = append(actions,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => false, });").Do(ctx)
			return err
		}),
		emulation.SetUserAgentOverride(userAgent).WithAcceptLanguage("pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7"),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			scrollAmount := 100 + rand.Intn(200)
			script := fmt.Sprintf("window.scrollBy(0, %d);", scrollAmount)
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		chromedp.Sleep(time.Duration(100+rand.Intn(100))*time.Millisecond),
		chromedp.Click("body", chromedp.ByQuery, chromedp.AtLeast(0)),
	)
	if options.SaveBody {
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error { return nil }))
	}
	return actions
}

func (b *Browser) resolveUserAgent(options types.FetchOptions) string {
	if options.UserAgent != "" {
		return options.UserAgent
	}
	if b.GetConfig().RandomizeUserAgent {
		return b.getRandomUserAgent()
	}
	return userAgents[0]
}

func (b *Browser) prepareContext(ctx context.Context, options types.FetchOptions) context.Context {
	ctx = context.WithValue(ctx, AllowRequestResources, options.AllowAssets)
	return ctx
}

func (b *Browser) GetCookies(ctx context.Context) ([]*network.Cookie, error) {
	if !b.initialized {
		b.logger.Error("browser not initialized")
		return nil, types.ErrInitBrowser
	}
	var cookies []*network.Cookie
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		b.logger.Error("error fetching cookies", "err", err)
		return nil, types.ErrGetCookies
	}
	return cookies, nil
}

func (b *Browser) Close() {
	if b.allocCancel != nil {
		b.allocCancel()
	}
	b.initialized = false
}
