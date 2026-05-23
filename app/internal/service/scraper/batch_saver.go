package scraper

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	factoryInterfaces "example-wikipedia-scraper/internal/interfaces/factory"
	repositoryInterace "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/db"
	"example-wikipedia-scraper/pkg/helpers"

	"slices"
)

type BatchSaver struct {
	logger       interfaces.LoggerInterface
	pageFactory  factoryInterfaces.PageFactoryInterface
	repository   repositoryInterace.PageRepositoryInterface
	failedPages  chan *dto.UnprocessedPageDTO
	queryBuilder db.QueryBuilder
}

func NewBatchSaver(logger interfaces.LoggerInterface, pageFactory factoryInterfaces.PageFactoryInterface, repository repositoryInterace.PageRepositoryInterface, failedPages chan *dto.UnprocessedPageDTO, queryBuilder db.QueryBuilder) *BatchSaver {
	return &BatchSaver{
		logger:       logger,
		pageFactory:  pageFactory,
		repository:   repository,
		failedPages:  failedPages,
		queryBuilder: queryBuilder,
	}
}

func (b *BatchSaver) ProcessBatch(batch []*dto.PageDTO) []*dto.PageDTO {
	if len(batch) > 0 {
		b.SaveBatch(batch)
		batch = batch[:0]
	}
	return batch
}

func (b *BatchSaver) SaveBatch(pages []*dto.PageDTO) {
	if len(pages) == 0 {
		return
	}
	b.logger.Debug("Saving batch of pages to database", "pages", pages)
	pageModels := make([]*model.Page, 0, len(pages))
	for _, pageDTO := range pages {
		pageModel, err := b.pageFactory.CreateFromDTO(pageDTO)
		if err != nil {
			b.logger.Error("Failed to convert page DTO to model", "pageDTO", pageDTO, "error", err)
			b.failedPages <- &dto.UnprocessedPageDTO{
				SiteName: pageDTO.SiteName,
				URL:      pageDTO.URL,
			}
			continue
		}
		pageModels = append(pageModels, pageModel)
	}
	b.setSamePageNotified(pageModels)
	for _, page := range pageModels {
		err := b.repository.Upsert(page, "url")
		if err != nil {
			b.failedPages <- &dto.UnprocessedPageDTO{
				SiteName: page.SiteName,
				URL:      page.URL,
			}
			b.logger.Error("Failed to upsert page", "page", page, "error", err)
		} else {
			b.logger.Debug("Successfully upserted page", "page", page)
		}
	}
}

func (b *BatchSaver) setSamePageNotified(pages []*model.Page) {
	hashKeys := make([]string, 0, len(pages))
	for _, page := range pages {
		hashKeys = append(hashKeys, page.HashKey)
	}
	hashKeys = slices.Compact(hashKeys)
	existingHashKeys := make([]string, 0, len(hashKeys))
	b.queryBuilder.Model(&model.Page{}).Select("hash_key").Where("notified = true AND hash_key IN ?", helpers.ToStringSlice(hashKeys)).Find(&existingHashKeys)
markNotified:
	for _, page := range pages {
		for _, existingKey := range existingHashKeys {
			if page.HashKey == existingKey {
				page.Notified = true
				continue markNotified
			}
		}
	}
}

func (b *BatchSaver) savingErrorLog(err error, pages []*model.Page) {
	if err != nil {
		pageModelsSummary := make([]model.Page, len(pages))
		for i, page := range pages {
			pageModelsSummary[i] = *page
		}
		b.logger.Error("Error saving batch", "error", err, "batch", pageModelsSummary)
	} else {
		b.logger.Debug("Successfully saved pages", "count", len(pages))
	}
}
