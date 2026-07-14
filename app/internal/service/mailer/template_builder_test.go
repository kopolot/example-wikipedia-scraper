package mailer

import (
	"strings"
	"testing"

	"example-wikipedia-scraper/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateBuilderVerificationEmail(t *testing.T) {
	builder := NewTemplateBuilder()

	mail, err := builder.VerificationEmail("user@example.com", "http://localhost:3000/dashboard/verify-email?token=abc")
	require.NoError(t, err)

	assert.Equal(t, []string{"user@example.com"}, mail.To)
	assert.Equal(t, "Verify your email address", mail.Subject)
	assert.Contains(t, mail.Body, "Confirm your email address")
	assert.Contains(t, mail.Body, "http://localhost:3000/dashboard/verify-email?token=abc")
	assert.Contains(t, mail.Body, "Wiki Scraper Newsletter")
}

func TestTemplateBuilderPageNotificationEmail(t *testing.T) {
	builder := NewTemplateBuilder()
	page := model.Page{
		Model:    model.Model{ID: 15},
		Title:    "Bitwa pod Grunwaldem",
		SiteName: "wikipedia.pl",
		URL:      "https://pl.wikipedia.org/wiki/Bitwa_pod_Grunwaldem",
	}

	mail, err := builder.PageNotificationEmail(
		"noreply@example.com",
		[]string{"one@example.com", "two@example.com"},
		"http://localhost:8080/dashboard/",
		page,
	)
	require.NoError(t, err)

	assert.Equal(t, []string{"noreply@example.com"}, mail.To)
	assert.Equal(t, []string{"one@example.com", "two@example.com"}, mail.Bcc)
	assert.Equal(t, "New Wikipedia page matching your filters", mail.Subject)
	assert.Contains(t, mail.Body, "New matching page available")
	assert.Contains(t, mail.Body, "Bitwa pod Grunwaldem")
	assert.Contains(t, mail.Body, "wikipedia.pl")
	assert.Contains(t, mail.Body, "http://localhost:8080/dashboard/panel/page-records/15")
	assert.NotContains(t, mail.Body, "https://pl.wikipedia.org/wiki/Bitwa_pod_Grunwaldem")
}

func TestJoinURL(t *testing.T) {
	got, err := joinURL("http://localhost:8080/dashboard/", "/panel/page-records/", "/22/")
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8080/dashboard/panel/page-records/22", got)
	assert.False(t, strings.Contains(got, "//panel"))
}
