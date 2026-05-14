package mailer

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"

	types "example-wikipedia-scraper/internal/types/mail"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	mailerConfig *config.MailerConfig
	dialer       *gomail.Dialer
	dialerClose  gomail.SendCloser
	logger       interfaces.LoggerInterface
}

func NewMailer(mailerConfig *config.MailerConfig, logger interfaces.LoggerInterface) *Mailer {
	mailer := &Mailer{
		mailerConfig: mailerConfig,
		logger:       logger,
	}
	mailer.dialer = gomail.NewDialer(mailer.mailerConfig.SMTPHost, mailer.mailerConfig.SMTPPort, mailer.mailerConfig.Username, mailer.mailerConfig.Password)
	close, err := mailer.dialer.Dial()
	if err != nil {
		mailer.logger.Fatal("Failed to connect to mail server: " + err.Error())
	}
	mailer.dialerClose = close
	return mailer
}

func (m *Mailer) GetConfig() config.MailerConfig {
	return *m.mailerConfig
}

func (m *Mailer) NewMail() *types.Mail {
	return &types.Mail{}
}

func (m *Mailer) Send(mail *types.Mail) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.mailerConfig.SenderName+" <"+m.mailerConfig.SenderEmail+">")
	msg.SetHeader("To", mail.To...)
	if len(mail.Cc) > 0 {
		msg.SetHeader("Cc", mail.Cc...)
	}
	if len(mail.Bcc) > 0 {
		msg.SetHeader("Bcc", mail.Bcc...)
	}
	msg.SetHeader("Subject", mail.Subject)
	msg.SetBody("text/html", mail.Body)
	for _, file := range mail.AttachmentsFileName {
		msg.Attach(file)
	}
	return m.dialer.DialAndSend(msg)
}

func (m *Mailer) Close() error {
	return m.dialerClose.Close()
}
