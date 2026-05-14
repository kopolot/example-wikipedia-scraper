package mailer

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig = &config.MailerConfig{
	SMTPHost:    "mailhog",
	SMTPPort:    1025,
	Username:    "",
	Password:    "",
	SenderName:  "Test Sender",
	SenderEmail: "testsender@test.com",
}

func TestSend_Success(t *testing.T) {
	mailer := NewMailer(testConfig, &testutils.MockLogger{})
	mail := mailer.NewMail()
	mail.To = []string{"testrecipient@test.com"}
	mail.Subject = "Test Subject"
	mail.Body = "<h1>This is a test email</h1>"

	err := mailer.Send(mail)
	assert.NoError(t, err)
}
