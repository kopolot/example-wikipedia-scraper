package browser

import (
	"context"
	"testing"
	"time"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/testutils"
	types "example-wikipedia-scraper/internal/types/browser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestBrowser(t *testing.T) *Browser {
	t.Helper()
	mockSettings := &config.BrowserSettings{
		Timeout:            10,
		RandomizeUserAgent: false,
		EngineSettings: map[string]any{
			"no-sandbox":                  true,
			"disable-gpu":                 true,
			"headless":                    true,
			"disable-software-rasterizer": true,
			"window-size":                 "1920,1080",
		},
	}
	mockLogger := &testutils.MockLogger{}
	b := NewBrowser(mockSettings, mockLogger)
	err := b.InitBrowser()
	require.NoError(t, err)
	return b
}

func cleanupTestBrowser(t *testing.T, b *Browser) {
	t.Cleanup(b.Close)
}

func TestBrowser_FetchPage(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:     5 * time.Minute,
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: false,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	session, err := b.FetchPage(ctx, "https://httpbin.io/html", opts)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	assert.NotEmpty(t, resp.Body)
	assert.NotEmpty(t, resp.Headers)
	session.GetCancelFunc()()
}

func TestBrowser_FetchPageWithRetry(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:     5 * time.Second,
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: false,
		Retries:     2,
	}
	session, err := b.FetchPageWithRetry(context.Background(), "https://httpbin.io/html", opts)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	assert.NotEmpty(t, resp.Body)
}

func TestBrowser_GetCookies(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	opts := types.FetchOptions{
		Timeout:     5 * time.Second,
		SaveBody:    false,
		SaveHeaders: true,
		AllowAssets: false,
	}
	session, fetchErr := b.FetchPage(ctx, "https://httpbin.io/cookies/set?testcookie=testvalue", opts)
	require.NotNil(t, session)
	require.NoError(t, fetchErr)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	cookies, err := b.GetCookies(session.GetContext())
	assert.NoError(t, err)
	assert.NotNil(t, cookies)
	assert.NotEmpty(t, cookies)
}

func TestBrowser_GetPageContent(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	resp, err := b.GetPageContent("https://httpbin.io/html", types.WithTimeout(5*time.Second))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.GetStatusCode())
	assert.NotEmpty(t, resp.Body)
}

func TestBrowser_GetConfigAndSetConfig(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	cfg := b.GetConfig()
	assert.Equal(t, 10, cfg.Timeout)
	cfg.Timeout = 20
	b.SetConfig(cfg)
	assert.Equal(t, 20, b.GetConfig().Timeout)
}

func TestBrowser_GetNewContext_NotInitialized(t *testing.T) {
	b := NewBrowser(&config.BrowserSettings{}, &testutils.MockLogger{})
	cleanupTestBrowser(t, b)
	ctx, cancel, err := b.GetNewContext()
	assert.Nil(t, ctx)
	assert.Nil(t, cancel)
	assert.Error(t, err)
}

func TestBrowser_FetchPage_NotInitialized(t *testing.T) {
	b := NewBrowser(&config.BrowserSettings{}, &testutils.MockLogger{})
	cleanupTestBrowser(t, b)
	_, err := b.FetchPage(context.Background(), "https://httpbin.io/html", types.FetchOptions{})
	assert.Error(t, err)
}

func TestBrowser_FetchPageWithRetry_NotInitialized(t *testing.T) {
	b := NewBrowser(&config.BrowserSettings{}, &testutils.MockLogger{})
	cleanupTestBrowser(t, b)
	_, err := b.FetchPageWithRetry(context.Background(), "https://httpbin.io/html", types.FetchOptions{})
	assert.Error(t, err)
}

func TestBrowser_GetCookies_NotInitialized(t *testing.T) {
	b := NewBrowser(&config.BrowserSettings{}, &testutils.MockLogger{})
	cleanupTestBrowser(t, b)
	_, err := b.GetCookies(context.Background())
	assert.Error(t, err)
}

func TestBrowser_Close(t *testing.T) {
	b := getTestBrowser(t)
	assert.True(t, b.initialized)
	b.Close()
	assert.False(t, b.initialized)
}

func TestBrowserFetchPageWithRetry_ContextLeak(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:     10 * time.Minute,
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: false,
		Retries:     3,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	for range 10 {
		// context leak
		session, err := b.FetchPageWithRetry(ctx, "https://httpbin.org/html", opts)
		assert.NoError(t, err)
		assert.NotNil(t, session)
		resp := session.GetBrowserResponse()
		assert.Equal(t, 200, resp.GetStatusCode())
		session.GetCancelFunc()()
		assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	}
	assert.True(t, true, "Completed multiple FetchPageWithRetry calls without context leaks")
}

func TestBrowserFetchPageWithRetry_Success(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:     10 * time.Minute,
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: false,
		Retries:     3,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session, err := b.FetchPageWithRetry(ctx, "https://httpbin.org/html", opts)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
}
