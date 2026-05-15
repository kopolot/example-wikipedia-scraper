package config

import (
	"encoding/json"
	"example-wikipedia-scraper/pkg/helpers"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigInterface interface {
	GetConfigFile() string
	DecodeFile(file *os.File) error
	GetDBString() string
	GetSitesConfig() []*SiteConfig
	GetBrowserSettings() *BrowserSettings
	IsDBDebugMode() bool
	GetScraperSettings() *ScraperSettings
	GetApiConfig() *ApiConfig
	GetMailerConfig() *MailerConfig
	GetNotifierConfig() *NotifierConfig
	GetDBConfig() *DBConfig
	GetRabbitMQConfig() *RabbitMQConfig
}

func (c *BrowserSettings) GetBrowserEngineSettings() map[string]any {
	if c.EngineSettings == nil {
		c.EngineSettings = make(map[string]any)
	}
	return c.EngineSettings
}

type Config struct {
	SitesToScrape   []*SiteConfig    `json:"sites_to_scrape"`
	DB              *DBConfig        `json:"db"`
	BrowserSettings *BrowserSettings `json:"browser_settings"`
	ScraperSettings *ScraperSettings `json:"scraper_settings"`
	ApiConfig       *ApiConfig       `json:"api"`
	MailerConfig    *MailerConfig    `json:"mailer"`
	NotifierConfig  *NotifierConfig  `json:"notifier"`
	RabbitMQConfig  *RabbitMQConfig  `json:"rabbitmq"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	filePath := config.GetConfigFile()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := config.DecodeFile(file); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetConfigFile() string {

	currentDir := helpers.GetCurrentFilePath()
	return filepath.Join(currentDir, "..", "..", "config.json")
}

func (c *Config) DecodeFile(file *os.File) error {
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	return err
}

func (c *Config) GetDBString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name)
}

func (c *Config) GetSitesConfig() []*SiteConfig {
	return c.SitesToScrape
}

func (c *Config) GetBrowserSettings() *BrowserSettings {
	return c.BrowserSettings
}

func (c *Config) IsDBDebugMode() bool {
	return c.DB.Debug
}

func (c *Config) GetNotifierConfig() *NotifierConfig {
	if c.NotifierConfig == nil {
		c.NotifierConfig = &NotifierConfig{
			WorkerCount: 5,
		}
	}
	return c.NotifierConfig
}

func (c *Config) GetScraperSettings() *ScraperSettings {
	defaultSettings := &ScraperSettings{
		ValidateAndUpdatePagesInterval: 3600000, // 1 hour
		ValidateAndUpdatePagesCooldown: 300000,  // 5 minutes
	}
	if c.ScraperSettings == nil {
		return defaultSettings
	}
	if c.ScraperSettings.ValidateAndUpdatePagesInterval == 0 {
		c.ScraperSettings.ValidateAndUpdatePagesInterval = defaultSettings.ValidateAndUpdatePagesInterval
	}
	if c.ScraperSettings.ValidateAndUpdatePagesCooldown == 0 {
		c.ScraperSettings.ValidateAndUpdatePagesCooldown = defaultSettings.ValidateAndUpdatePagesCooldown
	}
	return c.ScraperSettings
}

func (c *Config) GetApiConfig() *ApiConfig {
	if c.ApiConfig == nil {
		c.ApiConfig = &ApiConfig{
			Debug: false,
		}
	}
	return c.ApiConfig
}

func (c *Config) GetMailerConfig() *MailerConfig {
	if c.MailerConfig == nil {
		c.MailerConfig = &MailerConfig{}
	}
	return c.MailerConfig
}

func (c *Config) GetDBConfig() *DBConfig {
	return c.DB
}

func (c *Config) GetRabbitMQConfig() *RabbitMQConfig {
	return c.RabbitMQConfig
}
