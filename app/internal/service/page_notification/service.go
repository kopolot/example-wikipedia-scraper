package page_notification

import (
	"encoding/json"
	"iter"
	"math"
	"slices"
	"sync"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	interfacesQueue "example-wikipedia-scraper/internal/interfaces/queue"
	repoInterfaces "example-wikipedia-scraper/internal/interfaces/repository"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	pageHandler "example-wikipedia-scraper/internal/queue/handlers/page"
	"example-wikipedia-scraper/internal/service/mailer"
	"example-wikipedia-scraper/pkg/db"
)

type PageNotificationService struct {
	pageFilterService   serviceInterfaces.PageFilterServiceInterface
	pageRepo            repoInterfaces.PageRepositoryInterface
	mailerService       interfaces.MailerInterface
	subscriptionService serviceInterfaces.SubscriptionServiceInterface
	logger              interfaces.LoggerInterface
	config              config.ConfigInterface
	queryBuilder        db.QueryBuilder
	queueService        interfacesQueue.MessageQueueServiceInterface
}

func NewPageNotificationService(
	pageFilterService serviceInterfaces.PageFilterServiceInterface,
	pageRepo repoInterfaces.PageRepositoryInterface,
	mailerService interfaces.MailerInterface,
	logger interfaces.LoggerInterface,
	config config.ConfigInterface,
	queryBuilder db.QueryBuilder,
	queueService interfacesQueue.MessageQueueServiceInterface,
	subscriptionService serviceInterfaces.SubscriptionServiceInterface,
) *PageNotificationService {
	return &PageNotificationService{
		pageFilterService:   pageFilterService,
		pageRepo:            pageRepo,
		mailerService:       mailerService,
		subscriptionService: subscriptionService,
		logger:              logger,
		config:              config,
		queryBuilder:        queryBuilder,
		queueService:        queueService,
	}
}

func (s *PageNotificationService) RegisterQueueHandlers() error {
	pageQueueHandler := pageHandler.NewPageQueueHandler(s, s.logger)
	s.queueService.RegisterHandler("page_notification", pageQueueHandler.HandleNotifyUsersAboutNewPages)
	return nil
}

func (s *PageNotificationService) EnqueuePagesNotifications() error {
	pages, err := s.GetNotNotifiedPages()
	if err != nil {
		return err
	}
	if len(pages) == 0 {
		s.logger.Info("no new pages to notify")
		return nil
	}
	chunks := s.getPageChunks(pages)
	var wg sync.WaitGroup
	for chunk := range chunks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.processChunkAndEnqueueTask(chunk)
		}()
	}
	wg.Wait()
	return nil
}

func (s *PageNotificationService) getPageChunks(pages []*model.Page) iter.Seq[[]*model.Page] {
	notifierConfig := s.config.GetNotifierConfig()
	workerCount := notifierConfig.WorkerCount
	if workerCount < 1 {
		workerCount = 1
	}
	chunkSize := int(math.Ceil(float64(len(pages)) / float64(workerCount)))
	return slices.Chunk(pages, chunkSize)
}

func (s *PageNotificationService) processChunkAndEnqueueTask(chunk []*model.Page) {
	for _, page := range chunk {
		users, err := s.findMatchingFiltersForPageAndGetUsersToNotify(*page)
		if err != nil {
			s.logger.Error("failed to find matching filters and get users to notify", "err", err, "pageID", page.ID)
			continue
		}
		taskData := pageHandler.NotifyUsersAboutNewPagesPayload{
			UsersEmail: func() []string {
				emails := make([]string, len(users))
				for i, user := range users {
					emails[i] = user.Email
				}
				return emails
			}(),
			PageID: page.ID,
		}
		payloadBytes, err := json.Marshal(taskData)
		if err != nil {
			s.logger.Error("failed to marshal task payload", "err", err, "pageID", page.ID)
			continue
		}
		payload := interfacesQueue.JSONString(payloadBytes)
		task := &interfacesQueue.Task{
			Type:    "page_notification",
			Payload: payload,
		}
		err = s.queueService.Publish(task)
		if err != nil {
			s.logger.Error("failed to publish task to queue", "err", err, "pageID", page.ID)
		}
		err = s.MarkPageAndSimilarAsNotified(page)
		if err != nil {
			s.logger.Error("failed to mark page as notified", "err", err, "pageID", page.ID)
		}
	}
}

func (s *PageNotificationService) findMatchingFiltersForPageAndGetUsersToNotify(page model.Page) ([]*model.User, error) {
	filters, err := s.GetMatchedPageFilters(page)
	if err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return nil, nil
	}
	return s.GetUsersToNotifyByFilters(filters)
}

func (s *PageNotificationService) NotifyUsersAboutPage(emails []string, pageID uint) error {
	page, err := s.pageRepo.GetByID(pageID)
	if err != nil {
		s.logger.Error("failed to get page by ID", "err", err, "pageID", pageID)
		return err
	}
	if len(emails) == 0 {
		return nil
	}
	err = s.SendNotificationToUsersAboutPage(emails, *page)
	if err != nil {
		s.logger.Error("failed to send notification to users", "err", err, "emails", emails, "pageID", page.ID)
		return err
	}
	return nil
}

func (s *PageNotificationService) GetNotNotifiedPages() ([]*model.Page, error) {
	var pages []*model.Page
	var pageIDs []uint
	err := s.queryBuilder.
		Select("MIN(id)").
		Table("pages").
		Where("notified = ?", false).
		Group("hash_key").Find(&pageIDs)
	if err != nil {
		return nil, err
	}
	if len(pageIDs) == 0 {
		return []*model.Page{}, nil
	}
	err = s.queryBuilder.
		Where("id IN (?)", pageIDs).
		Find(&pages)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (s *PageNotificationService) MarkPageAndSimilarAsNotified(page *model.Page) error {
	return s.queryBuilder.Table("pages").
		Where("hash_key = ?", page.HashKey).
		Update("notified", true)
}

func (s *PageNotificationService) GetMatchedPageFilters(page model.Page) ([]*model.UserWantedPagesFilter, error) {
	return s.pageFilterService.FindMatchingFiltersForPage(page)
}

func (s *PageNotificationService) GetUsersToNotifyByFilters(filters []*model.UserWantedPagesFilter) ([]*model.User, error) {
	userIds := make([]uint, 0, len(filters))
	for _, filter := range filters {
		userIds = append(userIds, filter.UserID)
	}
	userIds = slices.Compact(userIds)
	return s.subscriptionService.GetUsersWithActiveSubscription(userIds)
}

func (s *PageNotificationService) SendNotificationToUsersAboutPage(emails []string, page model.Page) error {
	mail, err := mailer.NewTemplateBuilder().PageNotificationEmail(
		s.mailerService.GetConfig().SenderEmail,
		emails,
		s.config.GetApiConfig().PublicFrontendHost,
		page,
	)
	if err != nil {
		return err
	}
	return s.mailerService.Send(mail)
}
