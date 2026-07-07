package scraper

import (
	"encoding/json"
	"errors"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	queue "example-wikipedia-scraper/internal/interfaces/queue"
	types "example-wikipedia-scraper/internal/types/scraper"
	"time"
)

type RetryFailedPagePayload struct {
	SiteName   string `json:"site_name"`
	Url        string `json:"url"`
	RetryCount int    `json:"retry_count"`
	MaxRetries int    `json:"max_retries"`
}

// FailedPageProcessor obsługuje przetwarzanie i ponowną próbę stron (SRP)
type FailedPageProcessor struct {
	browser    interfaces.BrowserInterface
	scraperMgr *ScraperManager
	pageQueue  chan *dto.PageDTO
	queueSvc   queue.MessageQueueServiceInterface
	logger     interfaces.LoggerInterface
}

func NewFailedPageProcessor(
	browser interfaces.BrowserInterface,
	scraperMgr *ScraperManager,
	pageQueue chan *dto.PageDTO,
	queueSvc queue.MessageQueueServiceInterface,
	logger interfaces.LoggerInterface,
) *FailedPageProcessor {
	return &FailedPageProcessor{
		browser:    browser,
		scraperMgr: scraperMgr,
		pageQueue:  pageQueue,
		queueSvc:   queueSvc,
		logger:     logger,
	}
}

// SaveFailedPage publikuje nieudaną stronę do RabbitMQ z initial backoff
func (fpp *FailedPageProcessor) SaveFailedPage(failedPage *dto.UnprocessedPageDTO) error {
	return fpp.SaveFailedPageWithRetry(failedPage, 0)
}

// SaveFailedPageWithRetry publikuje nieudaną stronę z backoff'em
func (fpp *FailedPageProcessor) SaveFailedPageWithRetry(failedPage *dto.UnprocessedPageDTO, retryCount int) error {
	payload := &RetryFailedPagePayload{
		SiteName:   failedPage.SiteName,
		Url:        failedPage.URL,
		RetryCount: retryCount,
		MaxRetries: 5,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fpp.logger.Error("Error marshaling payload", "err", err)
		return err
	}

	task := &queue.Task{
		Type:    "retry_failed_page",
		Payload: queue.JSONString(payloadJSON),
	}

	delay := time.Until(fpp.CalculateNextRetryTime(retryCount))
	if delay < 0 {
		delay = 0
	}

	if err := fpp.queueSvc.PublishWithDelay(task, delay); err != nil {
		fpp.logger.Error("Error publishing to queue", "url", failedPage.URL, "err", err)
		return err
	}

	fpp.logger.Info("Published failed page to queue with delay", "url", failedPage.URL, "retry_count", retryCount, "delay_sec", delay.Seconds())
	return nil
}

func (fpp *FailedPageProcessor) CalculateNextRetryTime(retryCount int) time.Time {
	delays := []time.Duration{
		10 * time.Second,
		15 * time.Minute,
		15 * time.Minute,
		15 * time.Minute,
		15 * time.Minute,
	}

	if retryCount >= len(delays) {
		retryCount = len(delays) - 1
	}

	return time.Now().Add(delays[retryCount])
}

// RegisterHandlers rejestruje handlery dla RabbitMQ
func (fpp *FailedPageProcessor) RegisterHandlers() {
	fpp.queueSvc.RegisterHandler("retry_failed_page", fpp.handleRetryFailedPage)
}

// handleRetryFailedPage obsługuje przetwarzanie pojedynczej strony z queue'u
func (fpp *FailedPageProcessor) handleRetryFailedPage(task *queue.Task) {
	var payload RetryFailedPagePayload
	if err := json.Unmarshal([]byte(task.Payload), &payload); err != nil {
		fpp.logger.Error("Error unmarshaling retry task payload", "err", err)
		return
	}

	siteName := payload.SiteName
	pageUrl := payload.Url
	if siteName == "" || pageUrl == "" {
		fpp.logger.Error("Invalid payload: missing site name or page URL", "payload", payload)
		return
	}

	scraper, exists := fpp.scraperMgr.Get(siteName)
	if !exists {
		fpp.logger.Error("Scraper not found, deleting retry", "site", siteName)
		return
	}

	pageDTO, err := scraper.ScrapePageData(pageUrl)
	if err != nil {
		// Obsłuż błędy
		if errors.Is(err, types.ErrRatelimit) {
			fpp.logger.Warn("Rate limit during retry", "url", pageUrl)
			return
		}
		if errors.Is(err, types.ErrRecordNotFound) {
			fpp.logger.Info("Failed page no longer valid, skipping", "url", pageUrl)
			return
		}
		// Inne błędy - retry z backoff'em
		fpp.logger.Error("Error retrying page, republishing with backoff", "url", pageUrl, "retry_count", payload.RetryCount, "err", err)
		if payload.RetryCount < payload.MaxRetries {
			if err := fpp.SaveFailedPageWithRetry(&dto.UnprocessedPageDTO{
				SiteName: siteName,
				URL:      pageUrl,
			}, payload.RetryCount+1); err != nil {
				fpp.logger.Error("Error republishing to queue", "url", pageUrl, "err", err)
				panic(err)
			}
		} else {
			fpp.logger.Error("Max retries exceeded", "url", pageUrl)
		}
		return
	}

	// Sukces - wysyłamy do przetwarzania
	fpp.pageQueue <- pageDTO
}
