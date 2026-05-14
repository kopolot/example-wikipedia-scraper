package testutils

import (
	"context"
	"example-wikipedia-scraper/internal/types/browser"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockBrowserSession) GetBrowserResponse() browser.BrowserResponse {
	args := m.Called()
	return args.Get(0).(browser.BrowserResponse)
}

func (m *MockBrowserSession) SetNetworkResponseAndClearBrowserResponse(response *network.Response) {
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

func (m *MockBrowserSession) SetOptions(saveBody, saveHeaders, allowAssets bool) {
	m.Called(saveBody, saveHeaders, allowAssets)
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
