package main

import (
	"log"
	"time"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/queue"
	"example-wikipedia-scraper/internal/service"
	"example-wikipedia-scraper/internal/service/mailer"
	"example-wikipedia-scraper/internal/service/page_notification"
	pkgRepo "example-wikipedia-scraper/pkg/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	isDebugMode := cfg.GetScraperSettings().Debug
	if isDebugMode {
		logger.Init("notify", logger.LevelDebug, true)
		logger.Info("Notify service running in debug mode")
	} else {
		logger.Init("notify", logger.LevelInfo, true)
	}
	if err := db.InitDB(cfg); err != nil {
		logger.Fatal("could not initialize database", "err", err)
	}
	defer db.CloseDB()
	initQueue(cfg)
	pageNotificationService := setUpPageNotificationService(cfg)
	pageNotificationService.RegisterQueueHandlers()
	queue.GetMessageQueueService().Start()
	for {
		start := time.Now()
		err := pageNotificationService.EnqueuePagesNotifications()
		if err != nil {
			logger.Error("error notifying users about new pages", "err", err)
		}
		elapsed := time.Since(start)
		sleepDuration := time.Second - elapsed
		if sleepDuration > 0 {
			logger.Info("sleeping before next notification cycle", "duration", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}

func setUpPageNotificationService(cfg *config.Config) *page_notification.PageNotificationService {
	userWantedFilterRepo := repository.NewUserWantedPagesFilterRepository()
	userRepo := repository.NewUserRepository()
	pageRepo := repository.NewPageRepository()
	pageFilterService := service.NewPageFilterService(userWantedFilterRepo, pageRepo)
	loggerInstance := logger.GetLogger()
	mailerService := mailer.NewMailer(cfg.GetMailerConfig(), loggerInstance)
	subLvlRepo := pkgRepo.NewGenericRepository[*model.SubscriptionLevel]()
	subLvlProdRepo := pkgRepo.NewGenericRepository[*model.SubscriptionLevelProduct]()
	subscriptionService := service.NewSubscriptionService(userWantedFilterRepo, subLvlRepo, subLvlProdRepo, userRepo)
	return page_notification.NewPageNotificationService(
		pageFilterService,
		pageRepo,
		mailerService,
		loggerInstance,
		cfg,
		db.GetQueryBuilder(),
		queue.GetMessageQueueService(),
		subscriptionService,
	)
}

func initQueue(cfg config.ConfigInterface) {
	service := queue.NewRabbitMQService(cfg.GetRabbitMQConfig())
	queue.InitMessageQueueService(service)
}
