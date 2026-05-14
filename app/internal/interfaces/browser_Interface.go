package interfaces

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	types "example-wikipedia-scraper/internal/types/browser"
	"time"

	"github.com/chromedp/cdproto/network"
)

type BrowserInterface interface {
	GetConfig() config.BrowserSettings
	SetConfig(config config.BrowserSettings)
	InitBrowser() error
	GetPageContent(url string, opts ...types.FetchOption) (*types.BrowserResponse, error)
	FetchPage(ctx context.Context, url string, options types.FetchOptions) (BrowserSessionInterface, error)
	GetNewContext() (context.Context, context.CancelFunc, error)
	GetCookies(context.Context) ([]*network.Cookie, error)
	Close()
	GetOptions(opts ...types.FetchOption) types.FetchOptions
	FetchPageWithRetry(ctx context.Context, url string, options types.FetchOptions) (BrowserSessionInterface, error)
}

type BrowserSessionInterface interface {
	GetContext() context.Context
	GetCancelFunc() context.CancelFunc
	GetBrowserResponse() types.BrowserResponse
	SetNetworkResponseAndClearBrowserResponse(response *network.Response)
	SetURL(url string)
	GetURL() string
	SetRequestID(id network.RequestID)
	GetRequestID() network.RequestID
	SetStartTime(t time.Time)
	SetEndTime(t time.Time)
	SetBody(body []byte)
	SetUserAgent(ua string)
	SetOptions(saveBody, saveHeaders, allowAssets bool)
	ShouldSaveBody() bool
	ShouldSaveHeaders() bool
	ShouldAllowAssets() bool
}
