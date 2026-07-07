package interfaces

import (
	"example-wikipedia-scraper/internal/config"
	types "example-wikipedia-scraper/internal/types/browser"
	"context"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type BrowserInterface interface {
	GetConfig() config.BrowserSettings
	SetConfig(config config.BrowserSettings)
	InitBrowser() error
	GetPageContent(url string, opts ...types.FetchOption) (*types.BrowserResponse, error)
	FetchPage(browserSession BrowserSessionInterface) error
	GetNewContext() (context.Context, context.CancelFunc, error)
	GetCookies(context.Context) ([]*network.Cookie, error)
	Close()
	GetOptions(opts ...types.FetchOption) types.FetchOptions
	FetchPageWithRetry(browserSession BrowserSessionInterface) error
	RunActions(ctx context.Context, actions ...chromedp.Action) error
}

type BrowserSessionInterface interface {
	GetContext() context.Context
	GetCancelFunc() context.CancelFunc
	GetBrowserResponse() types.BrowserResponse
	SetNetworkResponse(response *network.Response)
	SetURL(url string)
	GetURL() string
	SetRequestID(id network.RequestID)
	GetRequestID() network.RequestID
	SetStartTime(t time.Time)
	SetEndTime(t time.Time)
	SetBody(body []byte)
	SetUserAgent(ua string)
	SetOptions(options types.FetchOptions)
	GetOptions() types.FetchOptions
	ShouldSaveBody() bool
	ShouldSaveHeaders() bool
	ShouldAllowAssets() bool
}
