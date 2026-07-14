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
		{Name: "wikipedia.pl", Enabled: true},
		{Name: "example", Enabled: true, ProxyURL: "socks5://127.0.0.1:1080"},
		{Name: "wikipedia.en", Enabled: true, ProxyURL: "socks5://127.0.0.1:1080"},
	}

	pool := NewBrowserPool(settings, sites, &testutils.MockLogger{})
	wikipediaBrowser := pool.GetForSite(sites[0])
	proxiedExampleBrowser := pool.GetForSite(sites[1])
	proxiedEnBrowser := pool.GetForSite(sites[2])

	if wikipediaBrowser == proxiedExampleBrowser {
		t.Fatal("expected wikipedia.pl to use a different browser than proxied example site")
	}
	if proxiedExampleBrowser != proxiedEnBrowser {
		t.Fatal("expected proxied sites to share the same browser for the same proxy")
	}

	wikipediaImpl, ok := wikipediaBrowser.(*Browser)
	if !ok {
		t.Fatalf("unexpected browser type %T", wikipediaBrowser)
	}
	proxiedImpl, ok := proxiedExampleBrowser.(*Browser)
	if !ok {
		t.Fatalf("unexpected browser type %T", proxiedExampleBrowser)
	}
	if wikipediaImpl.proxyURL != "" {
		t.Fatalf("expected empty proxy for wikipedia.pl, got %q", wikipediaImpl.proxyURL)
	}
	if proxiedImpl.proxyURL != "socks5://127.0.0.1:1080" {
		t.Fatalf("expected proxied site proxy, got %q", proxiedImpl.proxyURL)
	}
}

func TestSingleBrowserProvider(t *testing.T) {
	settings := &config.BrowserSettings{}
	browserInstance := NewBrowser(settings, &testutils.MockLogger{})
	provider := &SingleBrowserProvider{Browser: browserInstance}

	if provider.GetForSite(&config.SiteConfig{Name: "wikipedia.pl", ProxyURL: "socks5://127.0.0.1:1080"}) != browserInstance {
		t.Fatal("expected single browser provider to always return the same browser")
	}
}
