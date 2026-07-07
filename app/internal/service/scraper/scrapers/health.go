package scrapers

import (
	"example-wikipedia-scraper/internal/interfaces"
	"time"
)

var siteHealth interfaces.SiteHealthInterface

func SetSiteHealth(health interfaces.SiteHealthInterface) {
	siteHealth = health
}

func waitBeforeScrapeAttempt(siteName string) {
	if siteHealth != nil {
		siteHealth.BeforeAttempt(siteName)
	}
}

func recordScrapeSuccess(siteName string) {
	if siteHealth != nil {
		siteHealth.OnSuccess(siteName)
	}
}

func recordScrapeFailure(siteName string, err error) time.Duration {
	if siteHealth == nil {
		return 0
	}
	return siteHealth.OnFailure(siteName, err)
}

func shouldLogScrapeError(siteName string) bool {
	if siteHealth == nil {
		return true
	}
	return siteHealth.ShouldLogError(siteName)
}
