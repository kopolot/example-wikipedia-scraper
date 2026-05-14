package interfaces

import (
	"example-wikipedia-scraper/internal/config"

	types "example-wikipedia-scraper/internal/types/mail"
)

type MailerInterface interface {
	GetConfig() config.MailerConfig
	NewMail() *types.Mail
	Send(mail *types.Mail) error
	Close() error
}
