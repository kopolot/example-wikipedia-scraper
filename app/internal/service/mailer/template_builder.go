package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path"
	"strings"

	types "example-wikipedia-scraper/internal/types/mail"

	"example-wikipedia-scraper/internal/model"
)

type TemplateBuilder struct{}

type templateData struct {
	Title       string
	Preview     string
	Headline    string
	Lead        string
	ButtonLabel string
	ButtonURL   string
	Details     []mailDetail
	Footer      string
}

type mailDetail struct {
	Label string
	Value string
}

func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{}
}

func (b *TemplateBuilder) VerificationEmail(recipientEmail, verificationLink string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Verify your email address", templateData{
		Title:       "Verify your email address",
		Preview:     "Confirm your email address to activate your account.",
		Headline:    "Confirm your email address",
		Lead:        "Thanks for creating an account. Confirm your email address to finish setup and start receiving matching Wikipedia pages.",
		ButtonLabel: "Verify email",
		ButtonURL:   verificationLink,
		Footer:      "If you did not create this account, you can safely ignore this message.",
	})
}

func (b *TemplateBuilder) CreateUserEmail(recipientEmail, verificationLink string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Welcome to Wiki Scraper Newsletter", templateData{
		Title:       "Welcome to Wiki Scraper Newsletter",
		Preview:     "Verify your email address to activate your new account.",
		Headline:    "Your account is ready",
		Lead:        "Welcome to Wiki Scraper Newsletter. Confirm your email address to activate your account and start receiving new matching pages.",
		ButtonLabel: "Verify email",
		ButtonURL:   verificationLink,
		Footer:      "If you did not create this account, you can safely ignore this message.",
	})
}

func (b *TemplateBuilder) EmailVerifiedEmail(recipientEmail string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Email address verified", templateData{
		Title:    "Email address verified",
		Preview:  "Your email address has been verified successfully.",
		Headline: "Email verified",
		Lead:     "Your email address has been verified successfully. You can now use your account without restrictions.",
		Footer:   "You can now sign in and manage your saved filters from the dashboard.",
	})
}

func (b *TemplateBuilder) PasswordResetEmail(recipientEmail, resetLink string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Password reset request", templateData{
		Title:       "Password reset request",
		Preview:     "Use the secure link below to reset your password.",
		Headline:    "Reset your password",
		Lead:        "We received a request to reset your password. Use the button below to choose a new one.",
		ButtonLabel: "Reset password",
		ButtonURL:   resetLink,
		Footer:      "If you did not request a password reset, you can ignore this message.",
	})
}

func (b *TemplateBuilder) PasswordChangedEmail(recipientEmail string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Password changed", templateData{
		Title:    "Password changed",
		Preview:  "Your password has been changed successfully.",
		Headline: "Password updated",
		Lead:     "Your password has been changed successfully. If this was not you, reset your password immediately.",
	})
}

func (b *TemplateBuilder) LogoutEmail(recipientEmail string) (*types.Mail, error) {
	return b.buildMail(recipientEmail, "Logout notification", templateData{
		Title:    "Logout notification",
		Preview:  "You have been logged out of your account.",
		Headline: "You have been logged out",
		Lead:     "This is a confirmation that your session has ended successfully.",
	})
}

func (b *TemplateBuilder) PageNotificationEmail(senderEmail string, recipientEmails []string, frontendHost string, page model.Page) (*types.Mail, error) {
	pageLink := page.URL
	if dashboardLink, err := joinURL(frontendHost, "panel/page-records", fmt.Sprintf("%d", page.ID)); err == nil {
		pageLink = dashboardLink
	}

	return b.buildMailWithRecipients([]string{senderEmail}, recipientEmails, "New Wikipedia page matching your filters", templateData{
		Title:       "New Wikipedia page matching your filters",
		Preview:     "A new page matching your saved filters is now available.",
		Headline:    "New matching page available",
		Lead:        "A new Wikipedia page matching your saved filters has just been added. Open the dashboard to review the details.",
		ButtonLabel: "Open page",
		ButtonURL:   pageLink,
		Details: []mailDetail{
			{Label: "Title", Value: page.Title},
			{Label: "Source", Value: page.SiteName},
		},
		Footer: "You are receiving this notification because you have matching saved filters.",
	})
}

func (b *TemplateBuilder) buildMail(recipientEmail, subject string, data templateData) (*types.Mail, error) {
	return b.buildMailWithRecipients([]string{recipientEmail}, nil, subject, data)
}

func (b *TemplateBuilder) buildMailWithRecipients(to, bcc []string, subject string, data templateData) (*types.Mail, error) {
	body, err := renderTemplate(data)
	if err != nil {
		return nil, err
	}

	return &types.Mail{
		To:      to,
		Bcc:     bcc,
		Subject: subject,
		Body:    body,
	}, nil
}

func renderTemplate(data templateData) (string, error) {
	tpl, err := template.New("mail").Parse(mailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func joinURL(baseURL string, segments ...string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	basePath := strings.TrimSuffix(parsed.Path, "/")
	allSegments := []string{basePath}
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		allSegments = append(allSegments, strings.Trim(segment, "/"))
	}
	parsed.Path = path.Join(allSegments...)
	if !strings.HasSuffix(parsed.Path, "/") && strings.HasSuffix(baseURL, "/") && len(segments) == 0 {
		parsed.Path += "/"
	}

	return parsed.String(), nil
}

const mailTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f7fb;font-family:Arial,sans-serif;color:#1f2937;">
  <div style="max-width:640px;margin:0 auto;padding:32px 16px;">
    <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">{{ .Preview }}</div>
    <div style="background:#ffffff;border-radius:16px;padding:32px;box-shadow:0 10px 30px rgba(15,23,42,0.08);">
      <p style="margin:0 0 12px;font-size:12px;letter-spacing:0.08em;text-transform:uppercase;color:#2563eb;font-weight:700;">Wiki Scraper Newsletter</p>
      <h1 style="margin:0 0 16px;font-size:28px;line-height:1.2;color:#111827;">{{ .Headline }}</h1>
      <p style="margin:0 0 24px;font-size:16px;line-height:1.6;color:#4b5563;">{{ .Lead }}</p>
      {{ if .ButtonURL }}
      <p style="margin:0 0 24px;">
        <a href="{{ .ButtonURL }}" style="display:inline-block;padding:14px 22px;background:#2563eb;color:#ffffff;text-decoration:none;border-radius:10px;font-weight:700;">{{ .ButtonLabel }}</a>
      </p>
      <p style="margin:0 0 24px;font-size:13px;line-height:1.5;color:#6b7280;">If the button does not work, copy and paste this link into your browser:<br><a href="{{ .ButtonURL }}" style="color:#2563eb;">{{ .ButtonURL }}</a></p>
      {{ end }}
      {{ if .Details }}
      <table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="border-collapse:collapse;margin:0 0 24px;">
        {{ range .Details }}
        <tr>
          <td style="padding:10px 0;border-top:1px solid #e5e7eb;font-size:14px;color:#6b7280;width:140px;">{{ .Label }}</td>
          <td style="padding:10px 0;border-top:1px solid #e5e7eb;font-size:14px;color:#111827;font-weight:600;">{{ .Value }}</td>
        </tr>
        {{ end }}
      </table>
      {{ end }}
      {{ if .Footer }}
      <p style="margin:0;font-size:13px;line-height:1.6;color:#6b7280;">{{ .Footer }}</p>
      {{ end }}
    </div>
  </div>
</body>
</html>`
