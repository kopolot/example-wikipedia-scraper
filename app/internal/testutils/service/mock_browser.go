package testutils

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	browsertype "example-wikipedia-scraper/internal/types/browser"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/mock"
)

type MockBrowser struct {
	mock.Mock
}

func (m *MockBrowser) GetConfig() config.BrowserSettings {
	args := m.Called()
	return args.Get(0).(config.BrowserSettings)
}

func (m *MockBrowser) SetConfig(cfg config.BrowserSettings) {
	m.Called(cfg)
}

func (m *MockBrowser) InitBrowser() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBrowser) GetPageContent(url string, opts ...browsertype.FetchOption) (*browsertype.BrowserResponse, error) {
	args := m.Called(url, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*browsertype.BrowserResponse), args.Error(1)
}

func (m *MockBrowser) FetchPage(browserSession interfaces.BrowserSessionInterface) error {
	args := m.Called(browserSession)
	return args.Error(0)
}

func (m *MockBrowser) GetNewContext() (context.Context, context.CancelFunc, error) {
	args := m.Called()
	return args.Get(0).(context.Context), args.Get(1).(context.CancelFunc), args.Error(2)
}

func (m *MockBrowser) GetCookies(ctx context.Context) ([]*network.Cookie, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*network.Cookie), args.Error(1)
}

func (m *MockBrowser) Close() {
	m.Called()
}

func (m *MockBrowser) GetOptions(opts ...browsertype.FetchOption) browsertype.FetchOptions {
	args := m.Called(opts)
	return args.Get(0).(browsertype.FetchOptions)
}

func (m *MockBrowser) FetchPageWithRetry(browserSession interfaces.BrowserSessionInterface) error {
	args := m.Called(browserSession)
	return args.Error(0)
}

func (m *MockBrowser) RunActions(ctx context.Context, actions ...chromedp.Action) error {
	args := m.Called(ctx, actions)
	return args.Error(0)
}

type MockBrowserSession struct {
	mock.Mock
}

func (m *MockBrowserSession) GetContext() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockBrowserSession) GetCancelFunc() context.CancelFunc {
	args := m.Called()
	return args.Get(0).(context.CancelFunc)
}

func (m *MockBrowserSession) GetBrowserResponse() browsertype.BrowserResponse {
	args := m.Called()
	return args.Get(0).(browsertype.BrowserResponse)
}

func (m *MockBrowserSession) SetNetworkResponse(response *network.Response) {
	m.Called(response)
}

func (m *MockBrowserSession) SetURL(url string) {
	m.Called(url)
}

func (m *MockBrowserSession) GetURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBrowserSession) SetRequestID(id network.RequestID) {
	m.Called(id)
}

func (m *MockBrowserSession) GetRequestID() network.RequestID {
	args := m.Called()
	return args.Get(0).(network.RequestID)
}

func (m *MockBrowserSession) SetStartTime(t time.Time) {
	m.Called(t)
}

func (m *MockBrowserSession) SetEndTime(t time.Time) {
	m.Called(t)
}

func (m *MockBrowserSession) SetBody(body []byte) {
	m.Called(body)
}

func (m *MockBrowserSession) SetUserAgent(ua string) {
	m.Called(ua)
}

func (m *MockBrowserSession) SetOptions(options browsertype.FetchOptions) {
	m.Called(options)
}

func (m *MockBrowserSession) GetOptions() browsertype.FetchOptions {
	args := m.Called()
	return args.Get(0).(browsertype.FetchOptions)
}

func (m *MockBrowserSession) ShouldSaveBody() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockBrowserSession) ShouldSaveHeaders() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockBrowserSession) ShouldAllowAssets() bool {
	args := m.Called()
	return args.Bool(0)
}
