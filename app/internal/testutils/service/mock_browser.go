package testutils

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	browsertype "example-wikipedia-scraper/internal/types/browser"

	"github.com/chromedp/cdproto/network"
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
	return args.Get(0).(*browsertype.BrowserResponse), args.Error(1)
}
func (m *MockBrowser) FetchPage(ctx context.Context, url string, options browsertype.FetchOptions) (interfaces.BrowserSessionInterface, error) {
	args := m.Called(ctx, url, options)
	return args.Get(0).(interfaces.BrowserSessionInterface), args.Error(1)
}
func (m *MockBrowser) GetNewContext() (context.Context, context.CancelFunc, error) {
	args := m.Called()
	return args.Get(0).(context.Context), args.Get(1).(context.CancelFunc), args.Error(2)
}
func (m *MockBrowser) GetCookies(ctx context.Context) ([]*network.Cookie, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*network.Cookie), args.Error(1)
}
func (m *MockBrowser) Close() {
	m.Called()
}
func (m *MockBrowser) GetOptions(opts ...browsertype.FetchOption) browsertype.FetchOptions {
	args := m.Called(opts)
	return args.Get(0).(browsertype.FetchOptions)
}
func (m *MockBrowser) GetBrowserSession(ctx context.Context, options browsertype.FetchOptions) interfaces.BrowserSessionInterface {
	args := m.Called(ctx, options)
	return args.Get(0).(interfaces.BrowserSessionInterface)
}
func (m *MockBrowser) GetRequestContext(ctx context.Context, options browsertype.FetchOptions) (context.Context, context.CancelFunc) {
	args := m.Called(ctx, options)
	return args.Get(0).(context.Context), args.Get(1).(context.CancelFunc)
}
func (m *MockBrowser) FetchPageWithRetry(ctx context.Context, url string, options browsertype.FetchOptions) (interfaces.BrowserSessionInterface, error) {
	args := m.Called(ctx, url, options)
	return args.Get(0).(interfaces.BrowserSessionInterface), args.Error(1)
}
