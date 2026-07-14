package registry

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/service/scraper/scrapers"
	"sync"
)

type ScraperFactory func(url string) interfaces.ScraperInterface

var (
	scraperRegistry map[string]ScraperFactory
	initOnce        sync.Once
)

func NewScraperRegistry(provider interfaces.BrowserProvider, cfg config.ConfigInterface, loggerInstance interfaces.LoggerInterface) map[string]ScraperFactory {
	initOnce.Do(func() {
		mappedSites := func(cfg config.ConfigInterface) map[string]*config.SiteConfig {
			sites := make(map[string]*config.SiteConfig)
			for _, siteConfig := range cfg.GetSitesConfig() {
				sites[siteConfig.Name] = siteConfig
			}
			return sites
		}(cfg)
		scraperRegistry = map[string]ScraperFactory{
			"example": func(url string) interfaces.ScraperInterface {
				return scrapers.NewExampleScraper(url)
			},
			"wikipedia.pl": func(url string) interfaces.ScraperInterface {
				site := mappedSites["wikipedia.pl"]
				return scrapers.NewWikipediaPLScraper(url, provider.GetForSite(site), site, loggerInstance)
			},
		}
	})
	return scraperRegistry
}

func GetScraperRegistry() map[string]ScraperFactory {
	return scraperRegistry
}
