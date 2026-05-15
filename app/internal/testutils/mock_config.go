package testutils

import (
	"example-wikipedia-scraper/internal/config"
	"os"

	"github.com/stretchr/testify/mock"
)

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetConfigFile() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockConfig) DecodeFile(file *os.File) error {
	args := m.Called(file)
	return args.Error(0)
}
func (m *MockConfig) GetDBString() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockConfig) GetSitesConfig() []*config.SiteConfig {
	args := m.Called()
	return args.Get(0).([]*config.SiteConfig)
}
func (m *MockConfig) GetBrowserSettings() *config.BrowserSettings {
	args := m.Called()
	return args.Get(0).(*config.BrowserSettings)
}
func (m *MockConfig) IsDBDebugMode() bool {
	args := m.Called()
	return args.Bool(0)
}
func (m *MockConfig) GetScraperSettings() *config.ScraperSettings {
	args := m.Called()
	return args.Get(0).(*config.ScraperSettings)
}
func (m *MockConfig) GetApiConfig() *config.ApiConfig {
	args := m.Called()
	return args.Get(0).(*config.ApiConfig)
}
func (m *MockConfig) GetMailerConfig() *config.MailerConfig {
	args := m.Called()
	return args.Get(0).(*config.MailerConfig)
}
func (m *MockConfig) GetNotifierConfig() *config.NotifierConfig {
	args := m.Called()
	return args.Get(0).(*config.NotifierConfig)
}
func (m *MockConfig) GetDBConfig() *config.DBConfig {
	args := m.Called()
	return args.Get(0).(*config.DBConfig)
}

func (m *MockConfig) GetRabbitMQConfig() *config.RabbitMQConfig {
	args := m.Called()
	return args.Get(0).(*config.RabbitMQConfig)
}
