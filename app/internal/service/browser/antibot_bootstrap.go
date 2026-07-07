package browser

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	antiBotDetectionWindow = 3 * time.Second
	antiBotMaxWait         = 25 * time.Second
	antiBotStableDuration  = 2 * time.Second
	antiBotPollInterval    = 400 * time.Millisecond
)

var akamaiCookieNames = map[string]struct{}{
	"_abck":   {},
	"ak_bmsc": {},
	"bm_sz":   {},
	"bm_s":    {},
	"bm_sc":   {},
}

func antiBotBootstrapAction() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return waitForAntiBotBootstrap(ctx)
	})
}

func waitForAntiBotBootstrap(ctx context.Context) error {
	if !detectAkamaiCookies(ctx) {
		deadline := time.Now().Add(antiBotDetectionWindow)
		for time.Now().Before(deadline) {
			if err := ctx.Err(); err != nil {
				return err
			}
			if detectAkamaiCookies(ctx) {
				break
			}
			time.Sleep(antiBotPollInterval)
		}
		if !detectAkamaiCookies(ctx) {
			return nil
		}
	}

	deadline := time.Now().Add(antiBotMaxWait)
	var lastURL string
	var stableSince time.Time

	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return err
		}

		var currentURL string
		if err := chromedp.Location(&currentURL).Do(ctx); err != nil {
			return err
		}

		if currentURL != lastURL {
			lastURL = currentURL
			stableSince = time.Time{}
			time.Sleep(antiBotPollInterval)
			continue
		}

		if stableSince.IsZero() {
			stableSince = time.Now()
		}

		if isAbckValid(ctx) && time.Since(stableSince) >= antiBotStableDuration {
			return nil
		}

		time.Sleep(antiBotPollInterval)
	}

	return nil
}

func detectAkamaiCookies(ctx context.Context) bool {
	cookies, err := network.GetCookies().Do(ctx)
	if err != nil {
		return false
	}
	for _, c := range cookies {
		if _, ok := akamaiCookieNames[c.Name]; ok {
			return true
		}
	}
	return false
}

func isAbckValid(ctx context.Context) bool {
	cookies, err := network.GetCookies().Do(ctx)
	if err != nil {
		return false
	}
	for _, c := range cookies {
		if c.Name != "_abck" {
			continue
		}
		return isAbckCookieValueValid(c.Value)
	}
	return false
}

func isAbckCookieValueValid(value string) bool {
	if len(value) < 100 {
		return false
	}
	return true
}
