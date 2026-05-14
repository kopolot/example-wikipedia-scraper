package main

import (
	"example-wikipedia-scraper/internal/api"
	"example-wikipedia-scraper/internal/auth"
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
	logger.Init("api", logger.LevelInfo, false)
	loggerInstance := logger.GetLogger()
	defer db.CloseDB()
	if err := db.InitDB(cfg); err != nil {
		loggerInstance.Fatal("could not initialize database: ", err)
	}
	// if err := db.AutoMigrate(); err != nil {
	// 	loggerInstance.Fatal("could not migrate database: ", err)
	// }
	authManager := getAuthManager(cfg)
	initQueue(cfg)
	apiInstance := api.NewApi(cfg, loggerInstance, authManager)
	apiInstance.LoadModules()
	apiInstance.SetupRoutes()
	queue.GetMessageQueueService().Start()
	apiInstance.Run()
}

func getAuthManager(cfg *config.Config) *auth.AuthManager {
	userRepo := repository.NewUserRepository()
	return auth.NewAuthManager(cfg.GetApiConfig(), userRepo)
}

func initQueue(cfg config.ConfigInterface) {
	// service := queue.NewFakeMessageQueueService(1000)
	service := queue.NewRabbitMQService(cfg.GetRabbitMQConfig())
	queue.InitMessageQueueService(service)
}
