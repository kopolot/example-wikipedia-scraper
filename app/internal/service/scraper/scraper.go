package scraper

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/pkg/db"

	"context"
	"example-wikipedia-scraper/internal/interfaces"
	factoryInterfaces "example-wikipedia-scraper/internal/interfaces/factory"
	queue "example-wikipedia-scraper/internal/interfaces/queue"
	repositoryInterace "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/registry"
	"example-wikipedia-scraper/internal/service/browser"
	types "example-wikipedia-scraper/internal/types/scraper"
	"sync"
	"time"
)

// ScraperService orkiestruje proces scrapowania (facade pattern)
type ScraperService struct {
	config              config.ConfigInterface
	browser             interfaces.BrowserInterface
	logger              interfaces.LoggerInterface
	scraperMgr          *ScraperManager
	queueProcessor      *QueueProcessor
	failedPageProcessor *FailedPageProcessor
	pageValidator       *PageValidator
	scrapingWg          sync.WaitGroup
}

func NewScraperService(
	cfg config.ConfigInterface,
	logger interfaces.LoggerInterface,
	pageRepo repositoryInterace.PageRepositoryInterface,
	pageFactory factoryInterfaces.PageFactoryInterface,
	queryBuilder db.QueryBuilder,
	queueSvc queue.MessageQueueServiceInterface,
) *ScraperService {
	failedPagesChan := make(chan *dto.UnprocessedPageDTO, 1000)
	pageQueue := make(chan *dto.PageDTO, 1000)
	batchSaver := NewBatchSaver(logger, pageFactory, pageRepo, failedPagesChan, queryBuilder)

	scraperMgr := NewScraperManager(cfg, logger)
	queueProcessor := NewQueueProcessor(pageQueue, batchSaver, 10, 30*time.Second, logger)
	failedPageProcessor := NewFailedPageProcessor(nil, scraperMgr, pageQueue, queueSvc, logger)
	pageValidator := NewPageValidator(scraperMgr, pageRepo, pageFactory, logger)

	return &ScraperService{
		config:              cfg,
		logger:              logger,
		scraperMgr:          scraperMgr,
		queueProcessor:      queueProcessor,
		failedPageProcessor: failedPageProcessor,
		pageValidator:       pageValidator,
	}
}

func (s *ScraperService) Init() error {
	s.browser = browser.NewBrowser(s.config.GetBrowserSettings(), s.logger)
	scraperRegistry := registry.NewScraperRegistry(s.browser, s.config, s.logger)

	if err := s.browser.InitBrowser(); err != nil {
		s.logger.Error("Error initializing browser", "err", err)
		return err
	}

	s.scraperMgr.RegisterScrapers(scraperRegistry)
	s.failedPageProcessor.browser = s.browser
	s.failedPageProcessor.RegisterHandlers()

	s.queueProcessor.Start()

	return nil
}

func (s *ScraperService) RunScrapers() {
	s.logger.Info("Starting scrapers...")
	for siteName := range s.scraperMgr.GetAll() {
		s.scrapingWg.Add(1)
		go func(name string) {
			defer s.scrapingWg.Done()
			s.runScraper(name)
		}(siteName)
	}

	s.scrapingWg.Wait()
	s.queueProcessor.Stop()
	s.logger.Info("All scrapers finished")
}

func (s *ScraperService) RunScrapersInContinuousLoop() {
	s.logger.Info("Starting scrapers in continuous loop...")
	scraperRegistry := registry.NewScraperRegistry(s.browser, s.config, s.logger)

	for siteName := range s.scraperMgr.GetAll() {
		go func(name string) {
			for {
				s.runScraper(name)
				s.scraperMgr.Reload(name, scraperRegistry)
			}
		}(siteName)
	}

	select {}
}

func (s *ScraperService) runScraper(siteName string, opts ...types.ScrapeOption) {
	s.logger.Info("Starting scraper", "sitename", siteName)
	defer s.logger.Info("Scraper finished", "sitename", siteName)

	scraper, exists := s.scraperMgr.Get(siteName)
	if !exists {
		s.logger.Warn("Scraper not found", "sitename", siteName)
		return
	}

	failedPagesChan := make(chan *dto.UnprocessedPageDTO, 1000)
	defer close(failedPagesChan)

	// Goroutine obsługująca nieudane oferty
	go func() {
		for failedPage := range failedPagesChan {
			s.failedPageProcessor.SaveFailedPage(failedPage)
		}
	}()

	channels := &types.ScrapeChannels{
		PageQueue:   s.queueProcessor.pageQueue,
		FailedPages: failedPagesChan,
	}

	scraper.ScrapeAsync(channels, opts...)
}

func (s *ScraperService) ValidateAndUpdatePages(ctx context.Context) {
	config := s.config.GetScraperSettings()
	interval := time.Duration(config.ValidateAndUpdatePagesInterval) * time.Millisecond
	cooldown := time.Duration(config.ValidateAndUpdatePagesCooldown) * time.Millisecond

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.validatePageBatch(cooldown)
		}
	}
}

func (s *ScraperService) validatePageBatch(cooldown time.Duration) {
	config := s.config.GetScraperSettings()
	interval := time.Duration(config.ValidateAndUpdatePagesInterval) * time.Millisecond

	repo := s.pageValidator.repository
	pages, err := repo.GetPagesUpdatedBefore(-interval)
	if err != nil {
		s.logger.Error("Error fetching pages for validation", "err", err)
		return
	}

	s.logger.Info("Validating pages", "count", len(pages))

	for _, page := range pages {
		if err := s.pageValidator.ValidateAndUpdate(page); err != nil {
			s.logger.Debug("Validation error handled", "pageID", page.ID, "err", err)
		}
		time.Sleep(cooldown)
	}

	s.logger.Info("Validation cycle complete")
}

func (s *ScraperService) Shutdown() {
	s.logger.Info("Shutting down scraper service")
	s.queueProcessor.Stop()
	if s.browser != nil {
		s.browser.Close()
	}
	s.logger.Info("Scraper service shutdown complete")
}
