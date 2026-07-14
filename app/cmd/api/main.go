package main

import (
	"example-wikipedia-scraper/internal/api"
	"example-wikipedia-scraper/internal/auth"
	"example-wikipedia-scraper/internal/cache"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/queue"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	logger.Init("api", logger.LevelInfo, true)
	loggerInstance := logger.GetLogger()
	defer db.CloseDB()
	if err := db.InitDB(cfg); err != nil {
		loggerInstance.Fatal("could not initialize database: ", err)
	}
	authManager := getAuthManager(cfg)
	initRedis(cfg)
	initQueue(cfg)
	apiInstance := api.NewApi(cfg, loggerInstance, authManager)
	api.NewContainer(cfg, loggerInstance).LoadModules(apiInstance)
	apiInstance.SetupRoutes()
	queue.GetMessageQueueService().Start()
	apiInstance.Run()
}

func getAuthManager(cfg *config.Config) *auth.AuthManager {
	userRepo := repository.NewUserRepository()
	return auth.NewAuthManager(cfg.GetApiConfig(), userRepo)
}

func initRedis(cfg config.ConfigInterface) {
	if err := cache.InitRedis(cfg.GetRedisConfig()); err != nil {
		log.Fatalf("could not initialize redis: %v", err)
	}
}

func initQueue(cfg config.ConfigInterface) {
	service := queue.NewRabbitMQService(cfg.GetRabbitMQConfig())
	queue.InitMessageQueueService(service)
}
