package scrapers

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/service/browser"
	tested "example-wikipedia-scraper/internal/service/scraper/scrapers"
	scraperTypes "example-wikipedia-scraper/internal/types/scraper"
	"example-wikipedia-scraper/test/integration"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	wikiConfig      *config.SiteConfig
	browserInstance *browser.Browser
)

func setupWikipediaPLScraperTestDep(t *testing.T, cfg *config.Config) {
	t.Helper()
	browserInstance = browser.NewBrowser(cfg.GetBrowserSettings(), logger.GetLogger())
	sitesConfig := cfg.GetSitesConfig()
	for _, siteConfig := range sitesConfig {
		if siteConfig.Name == "wikipedia.pl" {
			wikiConfig = siteConfig
			break
		}
	}
	require.NotNil(t, wikiConfig, "wikipedia.pl config should not be nil")
}

func TestScrapeList_ExploreWillPageOffsetWork(t *testing.T) {
	logger.Init("test", logger.LevelDebug, true)
	cfg, err := integration.GetConfig()
	require.NoError(t, err)

	setupWikipediaPLScraperTestDep(t, cfg)

	err = browserInstance.InitBrowser()
	require.NoError(t, err)

	scraper := tested.NewWikipediaPLScraper(wikiConfig.URL, browserInstance, wikiConfig, logger.GetLogger())
	options := []scraperTypes.ScrapeOption{
		scraperTypes.WithMaxItems(1),
	}

	scraper.InitListScraperWorker(options...)
	err = scraper.StartListScraperWorker(options...)
	assert.NoError(t, err)
}

func TestWikipediaPLScraper_ScrapePage(t *testing.T) {

	logger.Init("test", logger.LevelDebug, true)
	cfg, err := integration.GetConfig()
	require.NoError(t, err)

	setupWikipediaPLScraperTestDep(t, cfg)

	err = browserInstance.InitBrowser()
	require.NoError(t, err)

	scraper := tested.NewWikipediaPLScraper(wikiConfig.URL, browserInstance, wikiConfig, logger.GetLogger())

	pageUrl := "https://pl.wikipedia.org/wiki/Chor%C4%85giew_pancerna_koronna_Stanis%C5%82awa_Wadowskiego"
	ctx, cancel, err := browserInstance.GetNewContext()
	require.NoError(t, err)
	defer cancel()
	dto, err := scraper.ScrapePageData(pageUrl, ctx)
	require.NoError(t, err)
	require.NotNil(t, dto)
}
