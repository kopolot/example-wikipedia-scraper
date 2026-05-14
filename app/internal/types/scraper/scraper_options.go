package scraper

import (
	"time"

	"github.com/chromedp/cdproto/network"
)

type ScrapeOption func(*ScrapeOptions)

type ScrapeOptions struct {
	MaxItems int
	Timeout  time.Duration
	Cookies  []*network.CookieParam
}

func WithMaxItems(max int) ScrapeOption {
	return func(o *ScrapeOptions) {
		o.MaxItems = max
	}
}

func WithTimeout(timeout time.Duration) ScrapeOption {
	return func(o *ScrapeOptions) {
		o.Timeout = timeout
	}
}

func WithCookies(cookies ...*network.CookieParam) ScrapeOption {
	return func(o *ScrapeOptions) {
		o.Cookies = cookies
	}
}

func ApplyOptions(opts ...ScrapeOption) *ScrapeOptions {
	options := &ScrapeOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
