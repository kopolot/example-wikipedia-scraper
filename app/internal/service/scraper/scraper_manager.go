package scraper

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/registry"
)

// ScraperManager zarządza zarejestrowaniem i przeładowaniem scraperów (SRP)
type ScraperManager struct {
	config     config.ConfigInterface
	scrapers   map[string]interfaces.ScraperInterface
	logger     interfaces.LoggerInterface
	siteHealth *SiteHealth
}

func NewScraperManager(cfg config.ConfigInterface, logger interfaces.LoggerInterface) *ScraperManager {
	return &ScraperManager{
		config:     cfg,
		scrapers:   make(map[string]interfaces.ScraperInterface),
		logger:     logger,
		siteHealth: NewSiteHealth(cfg, logger),
	}
}

func (sm *ScraperManager) GetSiteHealth() interfaces.SiteHealthInterface {
	return sm.siteHealth
}

func (sm *ScraperManager) RegisterScrapers(registry map[string]registry.ScraperFactory) {
	for _, siteConfig := range sm.config.GetSitesConfig() {
		if !siteConfig.Enabled {
			sm.logger.Info("Scraper for site is disabled, skipping", "sitename", siteConfig.Name)
			continue
		}

		factory, exists := registry[siteConfig.Name]
		if !exists {
			sm.logger.Info("No scraper registered for site", "sitename", siteConfig.Name)
			continue
		}

		sm.scrapers[siteConfig.Name] = factory(siteConfig.URL)
		sm.logger.Info("Registered scraper", "name", siteConfig.Name, "url", siteConfig.URL)
	}
}

func (sm *ScraperManager) Get(siteName string) (interfaces.ScraperInterface, bool) {
	scraper, exists := sm.scrapers[siteName]
	return scraper, exists
}

func (sm *ScraperManager) GetAll() map[string]interfaces.ScraperInterface {
	return sm.scrapers
}

func (sm *ScraperManager) Reload(siteName string, registry map[string]registry.ScraperFactory) {
	factory, exists := registry[siteName]
	if !exists {
		sm.logger.Info("No scraper registered for site, cannot reload", "sitename", siteName)
		return
	}

	siteConfig := sm.getSiteConfig(siteName)
	if siteConfig == nil {
		sm.logger.Info("Site configuration not found, cannot reload", "sitename", siteName)
		return
	}

	sm.scrapers[siteName] = factory(siteConfig.URL)
	sm.logger.Info("Reloaded scraper", "name", siteName)
}

func (sm *ScraperManager) getSiteConfig(siteName string) *config.SiteConfig {
	for _, siteConfig := range sm.config.GetSitesConfig() {
		if siteConfig.Name == siteName {
			return siteConfig
		}
	}
	return nil
}
