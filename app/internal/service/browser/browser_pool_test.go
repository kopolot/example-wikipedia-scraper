package browser

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/testutils"
	"testing"
)

func TestBrowserPoolUsesPerSiteProxy(t *testing.T) {
	settings := &config.BrowserSettings{
		EngineSettings: map[string]any{
			"headless": true,
		},
	}
	sites := []*config.SiteConfig{
		{Name: "otomoto", Enabled: true},
		{Name: "autoplac", Enabled: true, ProxyURL: "socks5://127.0.0.1:1080"},
		{Name: "olx", Enabled: true, ProxyURL: "socks5://127.0.0.1:1080"},
	}

	pool := NewBrowserPool(settings, sites, &testutils.MockLogger{})
	otomotoBrowser := pool.GetForSite(sites[0])
	autoplacBrowser := pool.GetForSite(sites[1])
	olxBrowser := pool.GetForSite(sites[2])

	if otomotoBrowser == autoplacBrowser {
		t.Fatal("expected otomoto to use a different browser than autoplac")
	}
	if autoplacBrowser != olxBrowser {
		t.Fatal("expected autoplac and olx to share the same browser for the same proxy")
	}

	otomotoImpl, ok := otomotoBrowser.(*Browser)
	if !ok {
		t.Fatalf("unexpected browser type %T", otomotoBrowser)
	}
	autoplacImpl, ok := autoplacBrowser.(*Browser)
	if !ok {
		t.Fatalf("unexpected browser type %T", autoplacBrowser)
	}
	if otomotoImpl.proxyURL != "" {
		t.Fatalf("expected empty proxy for otomoto, got %q", otomotoImpl.proxyURL)
	}
	if autoplacImpl.proxyURL != "socks5://127.0.0.1:1080" {
		t.Fatalf("expected autoplac proxy, got %q", autoplacImpl.proxyURL)
	}
}

func TestSingleBrowserProvider(t *testing.T) {
	settings := &config.BrowserSettings{}
	browserInstance := NewBrowser(settings, &testutils.MockLogger{})
	provider := &SingleBrowserProvider{Browser: browserInstance}

	if provider.GetForSite(&config.SiteConfig{Name: "autoplac", ProxyURL: "socks5://127.0.0.1:1080"}) != browserInstance {
		t.Fatal("expected single browser provider to always return the same browser")
	}
}
