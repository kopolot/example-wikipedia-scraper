package interfaces

import "time"

type SiteHealthInterface interface {
	BeforeAttempt(siteName string)
	OnFailure(siteName string, err error) time.Duration
	OnSuccess(siteName string)
	ShouldLogError(siteName string) bool
	IsCircuitOpen(siteName string) bool
}
