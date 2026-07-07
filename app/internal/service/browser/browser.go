package browser

// refractor i solid principles

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	types "example-wikipedia-scraper/internal/types/browser"
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ContextKey string

const (
	AllowRequestResources ContextKey = "allowRequestResources"
)

const defaultMaxConcurrentCDP = 8

type Browser struct {
	logger         interfaces.LoggerInterface
	allocCtx       context.Context
	chromeDpCtx    context.Context
	config         *config.BrowserSettings
	allocCancel    context.CancelFunc
	chromeDpCancel context.CancelFunc
	initMutex      sync.Mutex
	cdpSlots       chan struct{}
	initialized    bool
}

func NewBrowser(cfg *config.BrowserSettings, log interfaces.LoggerInterface) *Browser {
	return &Browser{
		config:   cfg,
		logger:   log,
		cdpSlots: make(chan struct{}, defaultMaxConcurrentCDP),
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
	b.initMutex.Lock()
	defer b.initMutex.Unlock()
	// ctx := context.Background()
	engineSettings := b.config.GetBrowserEngineSettings()
	b.applyDefaultEngineFlags(engineSettings)
	opts := make([]chromedp.ExecAllocatorOption, 0)
	// opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)
	// opts := []chromedp.ExecAllocatorOption{
	// 	chromedp.NoFirstRun,
	// 	chromedp.NoDefaultBrowserCheck,

	// 	// Stabilność
	// 	chromedp.Flag("disable-background-networking", true),
	// 	chromedp.Flag("disable-background-timer-throttling", true),
	// 	chromedp.Flag("disable-backgrounding-occluded-windows", true),
	// 	chromedp.Flag("disable-breakpad", true),
	// 	chromedp.Flag("disable-client-side-phishing-detection", true),
	// 	chromedp.Flag("disable-hang-monitor", true),
	// 	chromedp.Flag("disable-ipc-flooding-protection", true),
	// 	chromedp.Flag("disable-prompt-on-repost", true),
	// 	chromedp.Flag("disable-renderer-backgrounding", true),
	// 	chromedp.Flag("disable-sync", true),
	// 	chromedp.Flag("metrics-recording-only", true),
	// 	chromedp.Flag("safebrowsing-disable-auto-update", true),
	// 	chromedp.Flag("password-store", "basic"),
	// 	chromedp.Flag("use-mock-keychain", true),

	// 	// ↓ TO JEST KLUCZOWE – bez tych dwóch linii
	// 	// enable-automation NIE jest dodawane
	// 	chromedp.Flag("enable-automation", false),
	// 	chromedp.Flag("disable-blink-features", "AutomationControlled"),
	// 	chromedp.Flag("useAutomationExtension", false),
	// 	chromedp.Flag("disable-infobars", true),

	// 	// Twoje ustawienia
	// 	chromedp.Flag("headless", false),
	// 	chromedp.Flag("window-size", "1920,1080"),
	// 	chromedp.Flag("no-sandbox", true),
	// 	// chromedp.Flag("disable-setuid-sandbox", true),
	// 	chromedp.Flag("disable-dev-shm-usage", true),
	// 	chromedp.Flag("no-first-run", true),
	// 	chromedp.Flag("no-default-browser-check", true),
	// 	chromedp.Flag("disable-crash-reporter", true),

	// 	// GPU / WebGL
	// 	chromedp.Flag("use-gl", "swiftshader"),
	// 	chromedp.Flag("use-angle", "swiftshader-webgl"),
	// 	chromedp.Flag("enable-webgl", true),
	// 	chromedp.Flag("ignore-gpu-blocklist", true),

	// 	// CDP – losowy port żeby strony nie mogły go skanować
	// 	chromedp.Flag("remote-debugging-host", "127.0.0.1"),
	// 	chromedp.Flag("remote-debugging-port", getRandomPort()),
	// }

	cuOptions := []cu.Option{}
	for k, v := range engineSettings {
		if k == "headless" {
			cuOptions = append(cuOptions, cu.WithHeadless())
		} else {
			opts = append(opts, chromedp.Flag(k, v))
		}
	}
	cuOptions = append(cuOptions, cu.WithChromeFlags(opts...))

	ctx, cancel, err := cu.New(cu.NewConfig(cuOptions...))
	if err != nil {
		return err
	}
	b.chromeDpCtx = ctx
	b.chromeDpCancel = cancel
	// b.allocCtx, b.allocCancel = chromedp.NewExecAllocator(ctx, opts...)
	// chromedp.Run(b.allocCtx)
	// b.chromeDpCtx, b.chromeDpCancel = chromedp.NewContext(b.allocCtx)
	// to initialize the context*
	chromedp.Run(b.chromeDpCtx)
	b.initialized = true
	return nil
}

func getRandomPort() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err == nil {
		addr := l.Addr().String()
		l.Close() //nolint:errcheck,gosec

		return strings.Split(addr, ":")[1]
	}

	return "42069"
}

func (b *Browser) applyDefaultEngineFlags(settings map[string]any) {
	defaultFlags := map[string]any{
		"no-sandbox": true,
		// "disable-setuid-sandbox":   true,
		"disable-dev-shm-usage":    true,
		"disable-breakpad":         true,
		"disable-crash-reporter":   true,
		"no-first-run":             true,
		"no-default-browser-check": true,
	}
	for key, value := range defaultFlags {
		if _, exists := settings[key]; !exists {
			settings[key] = value
		}
	}
}

func (b *Browser) getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func (b *Browser) GetNewContext() (context.Context, context.CancelFunc, error) {
	if !b.initialized {
		return nil, nil, types.ErrInitBrowser
	}
	ctx := b.chromeDpCtx
	ctx, cancel := chromedp.NewContext(ctx)
	return ctx, cancel, nil
}

func (b *Browser) RunActions(ctx context.Context, actions ...chromedp.Action) error {
	b.acquireCDPSlot()
	defer b.releaseCDPSlot()
	return chromedp.Run(ctx, actions...)
}

func (b *Browser) acquireCDPSlot() {
	b.cdpSlots <- struct{}{}
}

func (b *Browser) releaseCDPSlot() {
	<-b.cdpSlots
}

func (b *Browser) shouldResetSessionContext(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "context canceled") ||
		strings.Contains(msg, "target closed") ||
		strings.Contains(msg, "cannot find context") ||
		strings.Contains(msg, "inspected target")
}

func (b *Browser) ResetSessionContext(browserSession interfaces.BrowserSessionInterface) error {
	session, ok := browserSession.(*BrowserSession)
	if !ok {
		return errors.New("invalid browser session type")
	}
	if cancel := session.GetCancelFunc(); cancel != nil {
		cancel()
	}
	newCtx, newCancel, err := b.GetNewContext()
	if err != nil {
		return err
	}
	session.ResetContext(newCtx, newCancel)
	return nil
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
	browserSession := NewBrowserSession(taskCtx, cancel)
	defer cancel()
	options := b.GetOptions(opts...)
	options.SaveBody = true

	browserSession.SetURL(url)
	browserSession.SetOptions(options)
	// browserSession := b.getBrowserSession(taskCtx, options)
	// requestCtx, requestCancel := b.getRequestContext(taskCtx, options)
	// defer requestCancel()
	response, err := b.runChromeDpWithActions(browserSession)
	if err != nil {
		b.logger.Error("error fetching page", "url", url, "err", err)
		return nil, err
	}
	b.logger.Debug("Fetched page content", "url", url, "status", response.GetStatusCode(), "responseTime", response.GetResponseTime())
	return response, nil
}

func (b *Browser) FetchPage(browserSession interfaces.BrowserSessionInterface) error {
	if !b.initialized {
		b.logger.Error("browser not initialized")
		return types.ErrInitBrowser
	}
	url := browserSession.GetURL()
	// browserSession := b.getBrowserSession(ctx, options)
	_, err := b.runChromeDpWithActions(browserSession)
	if err != nil {
		b.logger.Error("error fetching page", "url", url, "err", err)
		return err
	}
	b.logger.Debug("Fetched page content", "url", url)
	return nil
}

func (b *Browser) FetchPageWithRetry(browserSession interfaces.BrowserSessionInterface) error {
	if !b.initialized {
		return types.ErrInitBrowser
	}
	options := browserSession.GetOptions()
	url := browserSession.GetURL()
	var lastErr error
	var lastFetchErr error
	for attempt := range options.Retries {
		if attempt > 0 && b.shouldResetSessionContext(lastFetchErr) {
			if err := b.ResetSessionContext(browserSession); err != nil {
				b.logger.Warn("failed to reset browser session before retry", "url", url, "err", err)
			}
		}
		response, fetchErr := b.runChromeDpWithActions(browserSession)
		lastFetchErr = fetchErr
		err := b.handleFetchResult(response, url, fetchErr)
		switch {
		case errors.Is(err, types.ErrFetchPage):
			time.Sleep(1 * time.Second)
		case errors.Is(err, types.ErrFetchRatelimit):
			time.Sleep(b.getRetryCooldown(options))
		case errors.Is(err, types.ErrFetchManagedChallenge):
			b.logger.Warn("managed challenge detected, attempting to resolve", "url", url)
			solveErr := b.tryResolveManagedChallenge(browserSession, url, options)
			if solveErr == nil {
				response, fetchErr = b.runChromeDpWithActions(browserSession)
				lastFetchErr = fetchErr
				err = b.handleFetchResult(response, url, fetchErr)
				if err == nil {
					return nil
				}
			} else {
				b.logger.Warn("failed to resolve managed challenge", "url", url, "err", solveErr)
			}
		case errors.Is(err, types.ErrFetchTargetServer):
			time.Sleep(1 * time.Second)
		case errors.Is(err, nil):
			return nil
		case errors.Is(err, types.ErrFetchPageNotFound):
			return types.ErrFetchPageNotFound
		default:
			b.logger.Warn("unexpected error fetching page, retrying...", "url", url, "err", err)
			time.Sleep(1 * time.Second)
			err = types.ErrFetchPage
		}
		lastErr = err
	}
	return lastErr
}

func (b *Browser) cancelSession(browserSession interfaces.BrowserSessionInterface) {
	if browserSession == nil {
		return
	}
	cancel := browserSession.GetCancelFunc()
	if cancel != nil {
		cancel()
	}
}

func (b *Browser) getRetryCooldown(options types.FetchOptions) time.Duration {
	if options.RatelimitCooldown > 0 {
		return options.RatelimitCooldown
	}
	return 3 * time.Second
}

func (b *Browser) tryResolveManagedChallenge(browserSession interfaces.BrowserSessionInterface, url string, options types.FetchOptions) error {
	ctx := browserSession.GetContext()
	deadline := time.Now().Add(b.getChallengeResolveWindow(options))
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		// Na pierwszej iteracji daj Turnstile'owi czas na auto-solve przed próbą interakcji
		if attempt == 1 {
			time.Sleep(3 * time.Second)
		}

		didInteract, interactErr := b.performManagedChallengeInteraction(ctx)
		if interactErr != nil {
			b.logger.Warn("managed challenge interaction failed", "url", url, "attempt", attempt, "err", interactErr)
			break
		}

		waitDuration := time.Duration(1300+rand.Intn(900)) * time.Millisecond
		if didInteract {
			waitDuration = time.Duration(2100+rand.Intn(1300)) * time.Millisecond
			time.Sleep(waitDuration)
			return nil
		}
	}
	b.logger.Warn("managed challenge unresolved after deadline", "url", url)
	return types.ErrFetchManagedChallenge
}

func (b *Browser) getChallengeResolveWindow(options types.FetchOptions) time.Duration {
	const defaultWindow = 25 * time.Second
	if options.Timeout <= 0 {
		return defaultWindow
	}
	window := options.Timeout / 2
	if window < 12*time.Second {
		return 12 * time.Second
	}
	if window > 45*time.Second {
		return 45 * time.Second
	}
	return window
}

func (b *Browser) performManagedChallengeInteraction(ctx context.Context) (bool, error) {
	if err := b.RunActions(ctx, chromedp.WaitReady("body", chromedp.ByQuery)); err != nil {
		return false, err
	}

	b.RunActions(ctx,
		cu.MoveMouseToPosition(760+rand.Float64()*400, 440+rand.Float64()*200),
	)

	var iframe *cdp.Node

	err := b.RunActions(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		root, err := dom.GetDocument().
			WithDepth(-1).
			WithPierce(true).
			Do(ctx)
		if err != nil {
			return err
		}

		iframe = findIframe(root)
		if iframe == nil {
			return fmt.Errorf("iframe not found")
		}

		// fmt.Println("iframe node id:", iframe.NodeID)
		return nil
	}))

	if err != nil {
		panic(err)
	}

	err = b.RunActions(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		model, err := dom.GetBoxModel().
			WithNodeID(iframe.NodeID).
			Do(ctx)
		if err != nil {
			return err
		}
		x := (model.Content[0] + model.Content[4]) / 2
		y := (model.Content[1] + model.Content[5]) / 2
		err = input.DispatchMouseEvent(input.MousePressed, x, y).
			WithButton(input.Left).
			WithClickCount(1).
			Do(ctx)
		if err != nil {
			return err
		}

		err = input.DispatchMouseEvent(input.MouseReleased, x, y).
			WithButton(input.Left).
			WithClickCount(1).
			Do(ctx)
		return err
	}))

	if err != nil {
		return true, err
	}

	b.RunActions(ctx,
		cu.MoveMouseToPosition(760+rand.Float64()*400, 440+rand.Float64()*200),
	)
	return true, nil
}

func findIframe(n *cdp.Node) *cdp.Node {
	if n.NodeName == "IFRAME" {
		return n
	}

	for _, c := range n.Children {
		if r := findIframe(c); r != nil {
			return r
		}
	}

	for _, s := range n.ShadowRoots {
		if r := findIframe(s); r != nil {
			return r
		}
	}

	return nil
}

func (b *Browser) handleFetchResult(response *types.BrowserResponse, url string, err error) error {
	newErr := b.handleFetchError(response, url, err)
	if newErr != nil {
		return newErr
	}
	if b.isManagedChallenge(response) {
		b.logger.Warn("cloudflare managed challenge detected", "url", url, "status", response.GetStatusCode())
		return types.ErrFetchManagedChallenge
	}
	newErr = b.handleFetchStatusCode(response)
	if newErr != nil {
		return newErr
	}
	return nil
}

func (b *Browser) isManagedChallenge(response *types.BrowserResponse) bool {
	if response == nil {
		return false
	}
	if strings.Contains(strings.ToLower(response.URL), "/cdn-cgi/challenge-platform") {
		return true
	}
	body := strings.ToLower(response.GetBody())
	if body == "" {
		return false
	}
	// High-confidence: unique to Cloudflare, no false positives
	highConfidenceMarkers := []string{
		"challenges.cloudflare.com",
		"__cf_chl",
		"managed challenge",
		"cf-mitigated",
	}
	for _, marker := range highConfidenceMarkers {
		if strings.Contains(body, marker) {
			return true
		}
	}
	return false
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

// func (b *Browser) getBrowserSession(ctx context.Context, options types.FetchOptions) interfaces.BrowserSessionInterface {
// 	ctx, cancel := b.getRequestContext(ctx, options)
// 	browserSession := NewBrowserSession(ctx, cancel)
// 	return browserSession
// }

// tu trzeba szacher macer ultra mach bo cancel musi byc z chromedp, ale timeout trzeba dawac pozniej
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
	// ctx, secondCancel := chromedp.NewContext(ctx)
	// cancel := context.CancelFunc(func() {
	// 	firsrCancel()
	// 	secondCancel()
	// })
	return ctx, firsrCancel
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

func (b *Browser) runChromeDpWithActions(browserSession interfaces.BrowserSessionInterface) (*types.BrowserResponse, error) {
	options := browserSession.GetOptions()
	url := browserSession.GetURL()
	userAgent := b.resolveUserAgent(options)
	actionList := b.buildActionList(url, userAgent, options)
	ctx := browserSession.GetContext()
	session, ok := browserSession.(*BrowserSession)
	requestCtx, _ := b.getRequestContext(ctx, options)
	if !ok {
		// requestCancel()
		return nil, errors.New("invalid browser session type")
	}

	b.acquireCDPSlot()
	defer b.releaseCDPSlot()

	// session.SetUserAgent(userAgent)
	if options.BodyGetSleep > 0 {
		session.SetBodyGetSleep(options.BodyGetSleep)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	session.SetupResponseListener(&wg)
	response, err := chromedp.RunResponse(requestCtx, actionList...)
	// requestCancel()
	browserSession.SetNetworkResponse(response)
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
	// _, currentFile, _, _ := runtime.Caller(0)
	// currentDir := filepath.Dir(currentFile)
	// stealthScriptFileName := filepath.Join(helpers.FindRepoRoot(currentDir), "resource", "browser", "js", "stealth.js")
	// stealthScriptContent, err := os.ReadFile(stealthScriptFileName)
	// if err != nil {
	// 	b.logger.Warn("could not read stealth script, proceeding without it", "err", err)
	// }
	// stealthScript := string(stealthScriptContent)
	actions = append(actions,
		// to wywala ze mam headless
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	return network.SetExtraHTTPHeaders(network.Headers{
		// 		"accept-encoding": "gzip, deflate, br, zstd",
		// 		"accept-language": "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7",
		// 	}).Do(ctx)
		// }),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	return network.SetBlockedURLs([]string{
		// 		"*/127.0.0.1:*",
		// 		"*/localhost:*",
		// 	}).Do(ctx)
		// }),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, err := page.AddScriptToEvaluateOnNewDocument(stealthScript).Do(ctx)
		// 	return err
		// }),
		// emulation.SetUserAgentOverride(userAgent).
		// 	WithAcceptLanguage("pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7").
		// 	WithAcceptLanguage("pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7").
		// 	WithPlatform("Win32"),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.MouseEvent(
			input.MouseMoved,
			rand.Float64()*1920,
			rand.Float64()*1080,
		),
		chromedp.Sleep(time.Duration(100+rand.Intn(100))*time.Millisecond),
		chromedp.Click("body", chromedp.ByQuery, chromedp.AtLeast(0)),
	)
	if options.WaitAntiBotBootstrap {
		actions = append(actions, antiBotBootstrapAction())
	}
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
	err := b.RunActions(ctx,
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
	b.initMutex.Lock()
	defer b.initMutex.Unlock()
	if b.chromeDpCancel != nil {
		b.chromeDpCancel()
	}
	if b.allocCancel != nil {
		b.allocCancel()
	}
	b.initialized = false
}
