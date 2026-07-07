package browser

import (
	types "example-wikipedia-scraper/internal/types/browser"
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/stretchr/testify/assert"
)

func TestBrowserSession_BasicFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session := NewBrowserSession(ctx, cancel)
	session.SetURL("https://example.com")
	session.SetUserAgent("TestAgent")
	session.SetOptions(types.FetchOptions{
		SaveBody:    true,
		SaveHeaders: true,
		AllowAssets: false,
	})
	session.SetRequestID(network.RequestID("req-1"))
	session.SetStartTime(time.Now())
	session.SetEndTime(time.Now().Add(100 * time.Millisecond))
	session.SetBody([]byte("test body"))
	resp := &network.Response{
		Status:  200,
		Headers: map[string]interface{}{"Content-Type": "text/html"},
	}
	session.SetNetworkResponse(resp)
	br := session.GetBrowserResponse()
	assert.Equal(t, 200, br.StatusCode)
	assert.Equal(t, "https://example.com", br.URL)
	assert.Equal(t, "TestAgent", br.UserAgent)
	assert.Equal(t, "test body", br.Body)
	assert.Equal(t, "text/html", br.Headers["Content-Type"])
	assert.True(t, br.ResponseTime >= 100)
}

func TestBrowserSession_OptionsFlags(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session := NewBrowserSession(ctx, cancel)
	session.SetOptions(types.FetchOptions{
		SaveBody:    true,
		SaveHeaders: false,
		AllowAssets: true,
	})
	assert.True(t, session.ShouldSaveBody())
	assert.False(t, session.ShouldSaveHeaders())
	assert.True(t, session.ShouldAllowAssets())
}
