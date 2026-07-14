package interfaces

import "example-wikipedia-scraper/internal/config"

type BrowserProvider interface {
	GetForSite(site *config.SiteConfig) BrowserInterface
}
