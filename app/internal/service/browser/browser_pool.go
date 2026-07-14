package browser

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	"fmt"
)

type SingleBrowserProvider struct {
	Browser interfaces.BrowserInterface
}

func (p *SingleBrowserProvider) GetForSite(_ *config.SiteConfig) interfaces.BrowserInterface {
	return p.Browser
}

type BrowserPool struct {
	logger       interfaces.LoggerInterface
	baseSettings *config.BrowserSettings
	browsers     map[string]interfaces.BrowserInterface
	siteBrowsers map[string]interfaces.BrowserInterface
}

func NewBrowserPool(
	settings *config.BrowserSettings,
	sites []*config.SiteConfig,
	logger interfaces.LoggerInterface,
) *BrowserPool {
	pool := &BrowserPool{
		logger:       logger,
		baseSettings: settings,
		browsers:     make(map[string]interfaces.BrowserInterface),
		siteBrowsers: make(map[string]interfaces.BrowserInterface),
	}

	pool.ensureBrowser("")
	for _, site := range sites {
		if site == nil {
			continue
		}
		proxyKey := NormalizeProxyURL(site.ProxyURL)
		pool.ensureBrowser(proxyKey)
		pool.siteBrowsers[site.Name] = pool.browsers[proxyKey]
	}

	return pool
}

func (p *BrowserPool) ensureBrowser(proxyKey string) {
	if _, exists := p.browsers[proxyKey]; exists {
		return
	}
	p.browsers[proxyKey] = NewBrowser(p.baseSettings, p.logger, WithProxyURL(proxyKey))
}

func (p *BrowserPool) Init() error {
	for proxyKey, browserInstance := range p.browsers {
		if err := browserInstance.InitBrowser(); err != nil {
			if proxyKey == "" {
				return fmt.Errorf("init browser: %w", err)
			}
			return fmt.Errorf("init browser proxy=%q: %w", proxyKey, err)
		}
		if proxyKey != "" {
			p.logger.Info("Browser initialized with proxy", "proxy", proxyKey)
		}
	}
	return nil
}

func (p *BrowserPool) GetForSite(site *config.SiteConfig) interfaces.BrowserInterface {
	if site == nil {
		return p.browsers[""]
	}
	if browserInstance, ok := p.siteBrowsers[site.Name]; ok {
		return browserInstance
	}
	return p.browsers[NormalizeProxyURL(site.ProxyURL)]
}

func (p *BrowserPool) Close() {
	seen := make(map[interfaces.BrowserInterface]struct{}, len(p.browsers))
	for _, browserInstance := range p.browsers {
		if browserInstance == nil {
			continue
		}
		if _, exists := seen[browserInstance]; exists {
			continue
		}
		seen[browserInstance] = struct{}{}
		browserInstance.Close()
	}
}

var _ interfaces.BrowserProvider = (*BrowserPool)(nil)
var _ interfaces.BrowserProvider = (*SingleBrowserProvider)(nil)
