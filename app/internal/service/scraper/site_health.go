package scraper

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	types "example-wikipedia-scraper/internal/types/scraper"
	"errors"
	"sync"
	"time"
)

const (
	defaultFailureThreshold = 3
	defaultBaseBackoff      = 5 * time.Minute
	defaultMaxBackoff       = 30 * time.Minute
	defaultErrorLogInterval = 0
)

type SiteHealth struct {
	mu               sync.Mutex
	config           config.ConfigInterface
	logger           interfaces.LoggerInterface
	failures         map[string]int
	circuitOpenUntil map[string]time.Time
	lastErrorLogged  map[string]time.Time
	circuitOpenCount map[string]int
}

func NewSiteHealth(cfg config.ConfigInterface, logger interfaces.LoggerInterface) *SiteHealth {
	return &SiteHealth{
		config:           cfg,
		logger:           logger,
		failures:         make(map[string]int),
		circuitOpenUntil: make(map[string]time.Time),
		lastErrorLogged:  make(map[string]time.Time),
		circuitOpenCount: make(map[string]int),
	}
}

func (sh *SiteHealth) BeforeAttempt(siteName string) {
	for {
		wait, open := sh.waitDuration(siteName)
		if !open {
			return
		}
		if wait <= 0 {
			return
		}
		sh.logger.Info("Site unavailable, waiting before next scrape attempt",
			"sitename", siteName,
			"wait_sec", wait.Seconds(),
			"resume_at", sh.circuitOpenUntil[siteName].Format(time.RFC3339),
		)
		time.Sleep(wait)
	}
}

func (sh *SiteHealth) OnFailure(siteName string, err error) time.Duration {
	if !isSiteUnavailableError(err) {
		return 0
	}

	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.failures[siteName]++
	if sh.failures[siteName] < defaultFailureThreshold {
		return 0
	}

	backoff := sh.calculateBackoff(siteName)
	sh.circuitOpenCount[siteName]++
	sh.circuitOpenUntil[siteName] = time.Now().Add(backoff)
	sh.failures[siteName] = 0

	if sh.shouldLogErrorLocked(siteName) {
		sh.lastErrorLogged[siteName] = time.Now()
		sh.logger.Warn("Site marked as unavailable, backing off scrape attempts",
			"sitename", siteName,
			"err", err,
			"backoff_sec", backoff.Seconds(),
			"resume_at", sh.circuitOpenUntil[siteName].Format(time.RFC3339),
		)
	}

	return backoff
}

func (sh *SiteHealth) OnSuccess(siteName string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	wasOpen := sh.isCircuitOpenLocked(siteName)
	sh.failures[siteName] = 0
	sh.circuitOpenCount[siteName] = 0
	delete(sh.circuitOpenUntil, siteName)

	if wasOpen {
		sh.logger.Info("Site recovered, resuming scrape attempts", "sitename", siteName)
	}
}

func (sh *SiteHealth) ShouldLogError(siteName string) bool {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	return sh.shouldLogErrorLocked(siteName)
}

func (sh *SiteHealth) IsCircuitOpen(siteName string) bool {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	return sh.isCircuitOpenLocked(siteName)
}

func (sh *SiteHealth) waitDuration(siteName string) (time.Duration, bool) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	if !sh.isCircuitOpenLocked(siteName) {
		return 0, false
	}

	wait := time.Until(sh.circuitOpenUntil[siteName])
	if wait <= 0 {
		delete(sh.circuitOpenUntil, siteName)
		return 0, false
	}

	return wait, true
}

func (sh *SiteHealth) isCircuitOpenLocked(siteName string) bool {
	until, ok := sh.circuitOpenUntil[siteName]
	if !ok {
		return false
	}
	if time.Now().After(until) {
		delete(sh.circuitOpenUntil, siteName)
		return false
	}
	return true
}

func (sh *SiteHealth) shouldLogErrorLocked(siteName string) bool {
	last, ok := sh.lastErrorLogged[siteName]
	if !ok {
		return true
	}
	return time.Since(last) >= defaultErrorLogInterval
}

func (sh *SiteHealth) calculateBackoff(siteName string) time.Duration {
	base := defaultBaseBackoff
	for _, siteConfig := range sh.config.GetSitesConfig() {
		if siteConfig.Name == siteName && siteConfig.BlockedCooldown > 0 {
			base = time.Duration(siteConfig.BlockedCooldown) * time.Millisecond
			break
		}
	}

	multiplier := sh.circuitOpenCount[siteName]
	if multiplier > 4 {
		multiplier = 4
	}

	backoff := base * time.Duration(1<<multiplier)
	if backoff > defaultMaxBackoff {
		return defaultMaxBackoff
	}
	return backoff
}

func isSiteUnavailableError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, types.ErrRatelimit) ||
		errors.Is(err, types.ErrLastPageReached) ||
		errors.Is(err, types.ErrRecordNotFound) {
		return false
	}
	return errors.Is(err, types.ErrTargetServer) || errors.Is(err, types.ErrApplication)
}
