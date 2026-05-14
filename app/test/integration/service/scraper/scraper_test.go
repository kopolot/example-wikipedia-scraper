package scraper

import (
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/factory"
	"example-wikipedia-scraper/internal/queue"
	"example-wikipedia-scraper/internal/registry"

	"context"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model/repository"
	queuePkg "example-wikipedia-scraper/internal/queue"
	"example-wikipedia-scraper/internal/service/browser"
	"example-wikipedia-scraper/internal/service/scraper"
	"example-wikipedia-scraper/test/integration"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func handleGracefulShutdown(scraperService *scraper.ScraperService, cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logger.Info("Received shutdown signal, shutting down gracefully...")
		cancel()
		scraperService.Shutdown()
		close(done)
	}()
	<-done
}

func TestScraperService_InitializesSuccessfully(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),
		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queuePkg.NewFakeMessageQueueService(1000),
	)
	assert.NotNil(t, scraperService, "ScraperService should not be nil")
}

func TestInit_NoErrors(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),

		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queuePkg.NewFakeMessageQueueService(1000),
	)
	err = scraperService.Init()
	assert.NoError(t, err, "Expected no error during scraper service initialization")
}

func TestRunScrapers_CompletesSuccessfully(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),

		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queuePkg.NewFakeMessageQueueService(1000),
	)
	err = scraperService.Init()
	assert.NoError(t, err, "Expected no error during scraper service initialization")

	done := make(chan struct{})
	go func() {
		scraperService.RunScrapers()
		close(done)
	}()

	select {
	case <-done:
		t.Log("Scrapers completed successfully")
	case <-time.After(5 * time.Minute):
		t.Log("Scrapers timeout (5 minutes) - stopping")
		scraperService.Shutdown()
	}
}

func TestValidateAndUpdatePage_WithContext(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),

		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queuePkg.NewFakeMessageQueueService(1000),
	)
	err = scraperService.Init()
	assert.NoError(t, err, "Expected no error during scraper service initialization")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scraperService.ValidateAndUpdatePages(ctx)
	scraperService.Shutdown()
}

func TestContinuousScrapingWithValidation_NoErrors(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),
		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queuePkg.NewRabbitMQService(cfg.GetRabbitMQConfig()),
	)
	err = scraperService.Init()
	assert.NoError(t, err, "Expected no error during scraper service initialization")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	go scraperService.RunScrapersInContinuousLoop()
	go scraperService.ValidateAndUpdatePages(ctx)

	<-ctx.Done()
	scraperService.Shutdown()
}

func TestContinuousScraping_WithGracefulShutdown(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")
	scraperService := scraper.NewScraperService(
		cfg,
		logger.GetLogger(),

		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		// queuePkg.NewFakeMessageQueueService(1000),
		queue.NewRabbitMQService(cfg.GetRabbitMQConfig()),
	)
	err = scraperService.Init()
	assert.NoError(t, err, "Expected no error during scraper service initialization")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("Check logs for any errors in logs\n")
	go scraperService.RunScrapersInContinuousLoop()
	go scraperService.ValidateAndUpdatePages(ctx)

	handleGracefulShutdown(scraperService, cancel)
}

func TestAddFailedPageAndRetry(t *testing.T) {
	logger.Init("test", logger.LevelDebug, true)
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err, "Expected no error during test initialization")

	url := "https://example.com/failed-page"
	failedPage := &dto.UnprocessedPageDTO{
		SiteName: "example",
		URL:      url,
	}

	fmt.Printf("Adding failed page to queue: %v\n", failedPage)

	queueSvc := queue.NewRabbitMQService(cfg.GetRabbitMQConfig())
	browser := browser.NewBrowser(cfg.GetBrowserSettings(), logger.GetLogger())
	browser.InitBrowser()
	scraperMngr := scraper.NewScraperManager(cfg, logger.GetLogger())
	scraperMngr.RegisterScrapers(registry.NewScraperRegistry(browser, cfg, logger.GetLogger()))
	pageQueue := make(chan *dto.PageDTO, 1000)
	defer close(pageQueue)
	failedPageProc := scraper.NewFailedPageProcessor(browser, scraperMngr, pageQueue, queueSvc, logger.GetLogger())
	failedPageProc.RegisterHandlers()

	queueSvc.Start()
	go func() {
		for page := range pageQueue {
			fmt.Printf("Received page from failed page retry: %+v\n", page)
		}
	}()
	// time.Sleep(100 * time.Hour) // Daj chwilę na startowanie

	err = failedPageProc.SaveFailedPage(failedPage)

	time.Sleep(11 * time.Second)

	// page := <-pageQueue
	// fmt.Printf("Received page from failed page retry: %+v\n", page)

	err = failedPageProc.SaveFailedPage(&dto.UnprocessedPageDTO{
		SiteName: "aaaa",
		URL:      "aaaaa",
	})

	err = failedPageProc.SaveFailedPage(failedPage)

	time.Sleep(11 * time.Second)

	// page = <-pageQueue
	// fmt.Printf("Received page from failed page retry: %+v\n", page)

	time.Sleep(11 * time.Second)

	time.Sleep(10 * time.Minute)

	queueSvc.Close()
}
