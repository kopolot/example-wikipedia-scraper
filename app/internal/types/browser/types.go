package browser

import (
	"errors"
	"time"

	"github.com/chromedp/cdproto/network"
)

var (
	ErrFetchPage             = errors.New("browser error")
	ErrFetchRatelimit        = errors.New("ratelimited error")
	ErrFetchManagedChallenge = errors.New("managed challenge detected")
	ErrFetchTargetServer     = errors.New("target server error")
	ErrFetchPageNotFound     = errors.New("page not found")

	ErrGetPage        = errors.New("error getting page")
	ErrFetchWithRetry = errors.New("error fetching page with retry")
	ErrGetCookies     = errors.New("error getting cookies")

	ErrInitBrowser = errors.New("browser not initialized")
)

type FetchOptions struct {
	Cookies                []*network.CookieParam
	UserAgent              string
	Timeout                time.Duration
	RatelimitCooldown      time.Duration
	Retries                int
	SaveBody               bool
	SaveHeaders            bool
	AllowAssets            bool
	WaitAntiBotBootstrap   bool
	WarmupURL              string
	ConsentDismissScript   string
	BodyGetSleep           time.Duration
}

type FetchOption func(*FetchOptions)

func WithTimeout(timeout time.Duration) FetchOption {
	return func(o *FetchOptions) {
		o.Timeout = timeout
	}
}

func WithSaveBody(save bool) FetchOption {
	return func(o *FetchOptions) {
		o.SaveBody = save
	}
}

func WithSaveHeaders(save bool) FetchOption {
	return func(o *FetchOptions) {
		o.SaveHeaders = save
	}
}

func WithUserAgent(userAgent string) FetchOption {
	return func(o *FetchOptions) {
		o.UserAgent = userAgent
	}
}

func WithAllowAssets(allow bool) FetchOption {
	return func(o *FetchOptions) {
		o.AllowAssets = allow
	}
}

func WithRatelimitCooldown(cooldown time.Duration) FetchOption {
	return func(o *FetchOptions) {
		o.RatelimitCooldown = cooldown
	}
}

func WithCookies(cookies ...*network.CookieParam) FetchOption {
	return func(o *FetchOptions) {
		o.Cookies = cookies
	}
}

func WithRetries(retries int) FetchOption {
	return func(o *FetchOptions) {
		o.Retries = retries
	}
}

func WithBodyGetSleep(sleep time.Duration) FetchOption {
	return func(o *FetchOptions) {
		o.BodyGetSleep = sleep
	}
}

func WithWaitAntiBotBootstrap(wait bool) FetchOption {
	return func(o *FetchOptions) {
		o.WaitAntiBotBootstrap = wait
	}
}

func WithWarmupURL(url string) FetchOption {
	return func(o *FetchOptions) {
		o.WarmupURL = url
	}
}

func WithConsentDismissScript(script string) FetchOption {
	return func(o *FetchOptions) {
		o.ConsentDismissScript = script
	}
}

type BrowserResponse struct {
	UserAgent    string
	URL          string
	Body         string
	StatusCode   int
	Headers      map[string]string
	ResponseTime int64
	Response     *network.Response
}

func (r BrowserResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r BrowserResponse) GetBody() string {
	return r.Body
}

func (r BrowserResponse) GetHeaders() map[string]string {
	return r.Headers
}

func (r BrowserResponse) GetResponseTime() int64 {
	return r.ResponseTime
}

func (r BrowserResponse) GetURL() string {
	return r.URL
}

func (r BrowserResponse) GetUserAgent() string {
	return r.UserAgent
}

func (r BrowserResponse) GetResponse() network.Response {
	return *r.Response
}
