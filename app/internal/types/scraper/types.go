package scraper

import (
	"example-wikipedia-scraper/internal/dto"
)

type ScraperError string

func (e ScraperError) Error() string {
	return string(e)
}

const (
	ErrRecordNotFound  ScraperError = "record not found"
	ErrRatelimit       ScraperError = "ratelimited error"
	ErrApplication     ScraperError = "application error"
	ErrLastPageReached ScraperError = "last page reached"
	ErrTargetServer    ScraperError = "target server error"
)

type ScrapeChannels struct {
	PageQueue   chan<- *dto.PageDTO
	FailedPages chan<- *dto.UnprocessedPageDTO
}
