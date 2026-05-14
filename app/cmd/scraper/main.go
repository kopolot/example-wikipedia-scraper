package main

import (
	"context"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/factory"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/queue"
	"example-wikipedia-scraper/internal/service/scraper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	initLogger(cfg)

	if err := db.InitDB(cfg); err != nil {
		logger.Fatal("could not initialize database", "err", err)
	}
	defer db.CloseDB()

	initQueue(cfg)

	scraperService := createScraperService(cfg)
	queue.GetMessageQueueService().Start()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start all scraper components
	go scraperService.RunScrapersInContinuousLoop()
	go scraperService.ValidateAndUpdatePages(ctx)

	handleGracefulShutdown(scraperService, cancel)
}

func initLogger(cfg config.ConfigInterface) {
	isDebugMode := cfg.GetScraperSettings().Debug
	if isDebugMode {
		logger.Init("scraper", logger.LevelDebug, cfg.GetScraperSettings().CliLogging)
		logger.Info("Scraper running in debug mode")
	} else {
		logger.Init("scraper", logger.LevelInfo, cfg.GetScraperSettings().CliLogging)
	}
}

func createScraperService(cfg *config.Config) *scraper.ScraperService {
	loggerInstance := logger.GetLogger()
	queueHandler := queue.GetMessageQueueService()
	scraperService := scraper.NewScraperService(
		cfg,
		loggerInstance,
		repository.NewPageRepository(),
		factory.NewPageFactory(),
		db.GetQueryBuilder(),
		queueHandler,
	)
	err := scraperService.Init()
	if err != nil {
		loggerInstance.Fatal("could not create scraper service", "err", err)
	}
	return scraperService
}

func initQueue(cfg config.ConfigInterface) {
	service := queue.NewRabbitMQService(cfg.GetRabbitMQConfig())
	queue.InitMessageQueueService(service)
}

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
