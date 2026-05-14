package integration

import (
	"encoding/json"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/pkg/helpers"
	"os"
	"path/filepath"
)

type TestConfig struct {
	*config.Config
}

func (tc *TestConfig) GetConfigFile() string {
	currentDir := helpers.GetCurrentFilePath()
	return filepath.Join(currentDir, "config.test.json")
}

func (c *TestConfig) DecodeFile(file *os.File) error {
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	return err
}

func InitTest() (config.ConfigInterface, error) {
	logger.Init("test", logger.LevelDebug, true)
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	err = db.InitDB(cfg)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate()
	return cfg, nil
}

func GetConfig() (*config.Config, error) {
	cfg := &TestConfig{}
	filePath := cfg.GetConfigFile()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := cfg.DecodeFile(file); err != nil {
		return nil, err
	}
	return cfg.Config, nil
}

func CleanupDB() error {
	err := db.RollbackAllMigrations()
	if err != nil {
		logger.Error("could not rollback migrations", "err", err)
	}
	return err
}
