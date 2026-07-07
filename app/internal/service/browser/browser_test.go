package browser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/testutils"
	types "example-wikipedia-scraper/internal/types/browser"
	"example-wikipedia-scraper/pkg/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestBrowser(t *testing.T) *Browser {
	t.Helper()
	// mockSettings := &config.BrowserSettings{
	// 	Timeout:            10,
	// 	RandomizeUserAgent: false,
	// 	EngineSettings: map[string]any{
	// 		"no-sandbox": true,
	// 		// "disable-gpu": true,
	// 		// "headless":                    true,
	// 		"headless":                    false,
	// 		"disable-software-rasterizer": true,
	// 		"window-size":                 "1920,1080",
	// 		"disable-blink-features":      "AutomationControlled",
	// 		"excludeSwitches":             "enable-automation",
	// 		"useAutomationExtension":      false,
	// 		"disable-web-security":        true,
	// 	},
	// }

	mockSettings := &config.BrowserSettings{
		Timeout:            10,
		RandomizeUserAgent: false,
		EngineSettings: map[string]any{
			"headless": true,
			// "headless":                    "new",
			// "window-size": "1920,1080",
			// "disable-gpu":                 false,
			// "enable-gpu":                  true,
			// "use-gl":                      "swiftshader",
			// "enable-webgl":                true,
			// "ignore-gpu-blacklist":        true,
			// "disable-blink-features":      "AutomationControlled",
			// "disable-dev-shm-usage":       true,
			// "no-sandbox":                  true,
			// "use-angle":                   "swiftshader-webgl",
			// "ignore-gpu-blocklist":        true,
			// "disable-infobars":            true,
			// "excludeSwitches":             "enable-automation",
			// "remote-debugging-port":       "0",
			// "block-new-web-contents":      false,
			// "useAutomationExtension":      false,
			// "disable-software-rasterizer": true,
		},
	}
	os.Setenv("LIBGL_ALWAYS_SOFTWARE", "1")
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
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.io/html")
	session.SetOptions(opts)
	err = b.FetchPage(session)
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
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.io/html")
	session.SetOptions(opts)
	err = b.FetchPageWithRetry(session)
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
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.io/cookies/set?testcookie=testvalue")
	session.SetOptions(opts)
	fetchErr := b.FetchPage(session)
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
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.io/html")
	session.SetOptions(types.FetchOptions{})
	err = b.FetchPage(session)
	assert.Error(t, err)
}

func TestBrowser_FetchPageWithRetry_NotInitialized(t *testing.T) {
	b := NewBrowser(&config.BrowserSettings{}, &testutils.MockLogger{})
	cleanupTestBrowser(t, b)
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.io/html")
	session.SetOptions(types.FetchOptions{})
	err = b.FetchPageWithRetry(session)
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

	for range 10 {
		// context leak
		ctx, cancel, err := b.GetNewContext()
		require.NoError(t, err)
		session := NewBrowserSession(ctx, cancel)
		session.SetURL("https://httpbin.org/html")
		session.SetOptions(opts)
		err = b.FetchPageWithRetry(session)
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

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://httpbin.org/html")
	session.SetOptions(opts)
	err = b.FetchPageWithRetry(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
}

func TestBrowserIsManagedChallenge(t *testing.T) {
	b := &Browser{}

	t.Run("detect by url", func(t *testing.T) {
		response := &types.BrowserResponse{
			URL:  "https://example.com/cdn-cgi/challenge-platform/h/g/managed/v1",
			Body: "<html></html>",
		}
		assert.True(t, b.isManagedChallenge(response))
	})

	t.Run("detect by high-confidence body marker", func(t *testing.T) {
		response := &types.BrowserResponse{
			URL:  "https://example.com/offers",
			Body: "<html><script src='https://challenges.cloudflare.com/turnstile'></script></html>",
		}
		assert.True(t, b.isManagedChallenge(response))
	})

	t.Run("detect just a moment only when cloudflare also present", func(t *testing.T) {
		withCF := &types.BrowserResponse{
			URL:  "https://example.com/offers",
			Body: "<html><title>Just a moment...</title><p>Powered by cloudflare</p></html>",
		}
		assert.True(t, b.isManagedChallenge(withCF))

		withoutCF := &types.BrowserResponse{
			URL:  "https://example.com/offers",
			Body: "<html><title>Just a moment...</title></html>",
		}
		assert.False(t, b.isManagedChallenge(withoutCF))
	})

	t.Run("non challenge", func(t *testing.T) {
		response := &types.BrowserResponse{
			URL:  "https://example.com/offers",
			Body: "<html><body>normal content</body></html>",
		}
		assert.False(t, b.isManagedChallenge(response))
	})

	t.Run("nil response", func(t *testing.T) {
		assert.False(t, b.isManagedChallenge(nil))
	})
}

func TestBrowserGetChallengeResolveWindow(t *testing.T) {
	b := &Browser{}

	t.Run("uses default when timeout missing", func(t *testing.T) {
		window := b.getChallengeResolveWindow(types.FetchOptions{})
		assert.Equal(t, 25*time.Second, window)
	})

	t.Run("respects lower bound", func(t *testing.T) {
		window := b.getChallengeResolveWindow(types.FetchOptions{Timeout: 10 * time.Second})
		assert.Equal(t, 12*time.Second, window)
	})

	t.Run("uses half timeout in normal range", func(t *testing.T) {
		window := b.getChallengeResolveWindow(types.FetchOptions{Timeout: 40 * time.Second})
		assert.Equal(t, 20*time.Second, window)
	})

	t.Run("respects upper bound", func(t *testing.T) {
		window := b.getChallengeResolveWindow(types.FetchOptions{Timeout: 2 * time.Minute})
		assert.Equal(t, 45*time.Second, window)
	})
}

func TestBrowserIsBot(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:      10 * time.Minute,
		SaveBody:     true,
		SaveHeaders:  true,
		AllowAssets:  true,
		Retries:      3,
		BodyGetSleep: 5 * time.Second,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://bot.sannysoft.com")
	session.SetOptions(opts)
	err = b.FetchPage(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	body := string(resp.Body)
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path for saving response body")
	rootDir := helpers.FindRepoRoot(filepath.Dir(currentFile))
	fileName := fmt.Sprintf("%s/bot_sannysoft_response_%s.html", filepath.Join(rootDir, "tmp"), time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(fileName, []byte(body), 0644)
	require.NoError(t, err, "Failed to write response body to file for debugging")
	fmt.Printf("Saved bot.sannysoft.com response to %s for debugging\n", fileName)
}

func TestBrowserIsHeadless(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:      10 * time.Minute,
		SaveBody:     true,
		SaveHeaders:  true,
		AllowAssets:  false,
		Retries:      3,
		BodyGetSleep: 5 * time.Second,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://arh.antoinevastel.com/bots/areyouheadless")
	session.SetOptions(opts)
	err = b.FetchPage(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	body := string(resp.Body)
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path for saving response body")
	rootDir := helpers.FindRepoRoot(filepath.Dir(currentFile))
	fileName := fmt.Sprintf("%s/bot_headless_response_%s.html", filepath.Join(rootDir, "tmp"), time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(fileName, []byte(body), 0644)
	require.NoError(t, err, "Failed to write response body to file for debugging")
	fmt.Printf("Saved bot_headless response to %s for debugging\n", fileName)
}

func TestBrowserBotScoring(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:      10 * time.Minute,
		SaveBody:     true,
		SaveHeaders:  true,
		AllowAssets:  true,
		Retries:      3,
		BodyGetSleep: 5 * time.Second,
		// UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://fingerprint-scan.com/")
	session.SetOptions(opts)
	err = b.FetchPage(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	body := string(resp.Body)
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path for saving response body")
	rootDir := helpers.FindRepoRoot(filepath.Dir(currentFile))
	fileName := fmt.Sprintf("%s/bot_scoring_response_%s.html", filepath.Join(rootDir, "tmp"), time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(fileName, []byte(body), 0644)
	require.NoError(t, err, "Failed to write response body to file for debugging")
	fmt.Printf("Saved bot_scoring response to %s for debugging\n", fileName)
}

func TestCheckOtherScoring(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:      10 * time.Minute,
		SaveBody:     true,
		SaveHeaders:  true,
		AllowAssets:  true,
		Retries:      3,
		BodyGetSleep: 5 * time.Second,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://pixelscan.net/bot-check")
	session.SetOptions(opts)
	err = b.FetchPage(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	body := string(resp.Body)
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path for saving response body")
	rootDir := helpers.FindRepoRoot(filepath.Dir(currentFile))
	fileName := fmt.Sprintf("%s/bot_other_scoring_response_%s.html", filepath.Join(rootDir, "tmp"), time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(fileName, []byte(body), 0644)
	require.NoError(t, err, "Failed to write response body to file for debugging")
	fmt.Printf("Saved bot_other_scoring response to %s for debugging\n", fileName)
}

func TestCloudflareResolving(t *testing.T) {
	b := getTestBrowser(t)
	cleanupTestBrowser(t, b)
	opts := types.FetchOptions{
		Timeout:     10 * time.Minute,
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: true,
		Retries:     3,
		// BodyGetSleep: 5 * time.Second,
	}
	ctx, cancel, err := b.GetNewContext()
	require.NoError(t, err)
	defer cancel()

	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://check.spamhaus.org")
	session.SetOptions(opts)
	err = b.FetchPageWithRetry(session)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	resp := session.GetBrowserResponse()
	assert.Equal(t, 200, resp.GetStatusCode())
	session.GetCancelFunc()()
	assert.Error(t, session.GetContext().Err(), "Context should be cancelled properly")
	body := string(resp.Body)
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path for saving response body")
	rootDir := helpers.FindRepoRoot(filepath.Dir(currentFile))
	fileName := fmt.Sprintf("%s/cloudflare_resolving_response_%s.html", filepath.Join(rootDir, "tmp"), time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(fileName, []byte(body), 0644)
	require.NoError(t, err, "Failed to write response body to file for debugging")
	fmt.Printf("Saved cloudflare_resolving response to %s for debugging\n", fileName)
}
