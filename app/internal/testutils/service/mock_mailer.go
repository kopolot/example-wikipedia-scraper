package testutils

import (
	"example-wikipedia-scraper/internal/config"

	"github.com/stretchr/testify/mock"

	types "example-wikipedia-scraper/internal/types/mail"
)

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) GetConfig() config.MailerConfig {
	args := m.Called()
	return args.Get(0).(config.MailerConfig)
}

func (m *MockMailer) NewMail() *types.Mail {
	args := m.Called()
	return args.Get(0).(*types.Mail)
}

func (m *MockMailer) Send(mail *types.Mail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func (m *MockMailer) Close() error {
	args := m.Called()
	return args.Error(0)
}
