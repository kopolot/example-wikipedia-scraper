package page

import (
	"encoding/json"

	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/interfaces/queue"
	"example-wikipedia-scraper/internal/interfaces/service"
)

type NotifyUsersAboutNewPagesPayload struct {
	UsersEmail []string `json:"users_email" validate:"[]email"`
	PageID     uint     `json:"page_id" validate:"required"`
}

type PageQueueHandler struct {
	pageNotificationService service.PageNotificationServiceInterface
	logger                  interfaces.LoggerInterface
}

func NewPageQueueHandler(pageNotificationService service.PageNotificationServiceInterface, logger interfaces.LoggerInterface) *PageQueueHandler {
	return &PageQueueHandler{
		pageNotificationService: pageNotificationService,
		logger:                  logger,
	}
}

func (h *PageQueueHandler) HandleNotifyUsersAboutNewPages(task *queue.Task) {
	var payload NotifyUsersAboutNewPagesPayload
	err := json.Unmarshal([]byte(task.Payload), &payload)
	if err != nil {
		panic("Failed to unmarshal task payload: " + err.Error())
	}
	err = h.pageNotificationService.NotifyUsersAboutPage(payload.UsersEmail, payload.PageID)
	if err != nil {
		panic(err)
	}
}
