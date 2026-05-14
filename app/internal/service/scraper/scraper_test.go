package scraper

import (
	"testing"

	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	queue "example-wikipedia-scraper/internal/interfaces/queue"
	"example-wikipedia-scraper/internal/model"
	queueImpl "example-wikipedia-scraper/internal/queue"
	"example-wikipedia-scraper/internal/registry"
	"example-wikipedia-scraper/internal/testutils"

	testRepo "example-wikipedia-scraper/internal/testutils/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	cfg                *testutils.MockConfig
	logger             *testutils.MockLogger
	pageRepo           *testRepo.MockPageRepository
	pageFactoryScraper *testutils.MockPageFactory
	queryBuilder       *testutils.MockQueryBuilder
)

func setUpScraperService(t *testing.T) *ScraperService {
	t.Helper()
	cfg = &testutils.MockConfig{}
	logger = &testutils.MockLogger{}
	pageRepo = &testRepo.MockPageRepository{}
	pageFactoryScraper = &testutils.MockPageFactory{}
	queryBuilder = &testutils.MockQueryBuilder{}
	return NewScraperService(cfg, logger, pageRepo, pageFactoryScraper, queryBuilder, queueImpl.NewFakeMessageQueueService(1000))
}

func TestNewScraperService_InitializesComponents(t *testing.T) {
	s := setUpScraperService(t)
	assert.NotNil(t, s.scraperMgr)
	assert.NotNil(t, s.queueProcessor)
	assert.NotNil(t, s.failedPageProcessor)
	assert.NotNil(t, s.pageValidator)
}

func TestScraperManager_RegistersEnabledScrapers(t *testing.T) {
	mockScraper := &testutils.MockScraper{}
	registryInstance := map[string]registry.ScraperFactory{
		"site1": func(url string) interfaces.ScraperInterface { return mockScraper },
	}

	cfg := &testutils.MockConfig{}
	logger := &testutils.MockLogger{}
	cfg.On("GetSitesConfig").Return([]*config.SiteConfig{
		{Name: "site1", Enabled: true, URL: "http://a"},
		{Name: "site2", Enabled: false, URL: "http://b"},
	})

	scraperMgr := NewScraperManager(cfg, logger)
	scraperMgr.RegisterScrapers(registryInstance)

	scraper, exists := scraperMgr.Get("site1")
	assert.True(t, exists)
	assert.NotNil(t, scraper)

	_, notExists := scraperMgr.Get("site2")
	assert.False(t, notExists)
}

func TestValidateAndUpdatePage_Success(t *testing.T) {
	mockScraper := &testutils.MockScraper{}
	pageFactory := &testutils.MockPageFactory{}
	dto := &dto.PageDTO{ExternalID: "x", SiteName: "site"}
	page := &model.Page{ExternalID: "x", SiteName: "site"}
	pageFactory.On("CreateFromDTO", dto).Return(page, nil)
	mockScraper.On("ValidatePage", mock.Anything, mock.Anything).Return(dto, nil)
}

func TestPageComparator_DetectsChanges(t *testing.T) {
	comparator := NewPageComparator()

	page1 := &model.Page{ExternalID: "x", SiteName: "site"}
	page2 := &model.Page{ExternalID: "x", SiteName: "site"}
	page3 := &model.Page{ExternalID: "y", SiteName: "site"}

	assert.False(t, comparator.HasChanged(page1, page2))
	assert.True(t, comparator.HasChanged(page1, page3))
}

func TestFailedPageProcessor_SavesFailedPage(t *testing.T) {
	logger := &testutils.MockLogger{}
	queueSvc := &testutils.MockQueueService{}
	queueSvc.On("Publish", mock.Anything).Return(nil)

	processor := NewFailedPageProcessor(nil, nil, nil, queueSvc, logger)

	failedPage := &dto.UnprocessedPageDTO{
		SiteName: "test_site",
		URL:      "http://example.com/page",
	}

	err := processor.SaveFailedPage(failedPage)
	assert.NoError(t, err)
	queueSvc.AssertCalled(t, "Publish", mock.MatchedBy(func(task *queue.Task) bool {
		return task.Type == "retry_failed_page"
	}))
}
