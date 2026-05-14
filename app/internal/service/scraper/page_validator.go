package scraper

import (
	"errors"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	factoryInterfaces "example-wikipedia-scraper/internal/interfaces/factory"
	repositoryInterace "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/scraper"
	"time"
)

// PageValidator waliduje i aktualizuje strony (SRP)
type PageValidator struct {
	scraperMgr *ScraperManager
	repository repositoryInterace.PageRepositoryInterface
	comparator *PageComparator
	factory    factoryInterfaces.PageFactoryInterface
	logger     interfaces.LoggerInterface
}

func NewPageValidator(
	scraperMgr *ScraperManager,
	repository repositoryInterace.PageRepositoryInterface,
	factory factoryInterfaces.PageFactoryInterface,
	logger interfaces.LoggerInterface,
) *PageValidator {
	return &PageValidator{
		scraperMgr: scraperMgr,
		repository: repository,
		comparator: NewPageComparator(),
		factory:    factory,
		logger:     logger,
	}
}

func (pv *PageValidator) ValidateAndUpdate(page *model.Page) error {
	scraper, exists := pv.scraperMgr.Get(page.SiteName)
	if !exists {
		pv.logger.Warn("No scraper found for site, cannot validate page", "site", page.SiteName, "pageID", page.ID)
		return nil
	}

	dto, err := scraper.ValidatePage(page)
	if err != nil {
		return pv.handleValidationError(err, page)
	}

	if dto == nil {
		return nil
	}

	return pv.updatePageIfChanged(page, dto)
}

func (pv *PageValidator) handleValidationError(err error, page *model.Page) error {
	switch {
	case errors.Is(err, types.ErrRecordNotFound):
		return pv.deletePage(page)
	case errors.Is(err, types.ErrRatelimit):
		pv.logger.Warn("Rate limit during validation, stopping", "page", page.ID)
		return err
	default:
		pv.logger.Error("Error validating page", "pageID", page.ID, "err", err)
		return err
	}
}

func (pv *PageValidator) updatePageIfChanged(currentPage *model.Page, pageDTO *dto.PageDTO) error {
	newPage := pv.convertDTOToModel(currentPage, pageDTO)

	if !pv.comparator.HasChanged(currentPage, newPage) {
		return nil
	}

	if err := pv.repository.Update(newPage); err != nil {
		pv.logger.Error("Error updating page", "pageID", newPage.ID, "err", err)
		return err
	}

	pv.logger.Debug("Updated page with new data", "pageID", newPage.ID)
	return pv.markSimilarPagesAsUnnotified(currentPage)
}

func (pv *PageValidator) convertDTOToModel(currentPage *model.Page, pageDTO *dto.PageDTO) *model.Page {
	newPage, err := pv.factory.CreateFromDTO(pageDTO)
	if err != nil {
		pv.logger.Warn("Failed to convert DTO to model, keeping current data", "pageID", currentPage.ID, "err", err)
		return currentPage
	}

	newPage.ID = currentPage.ID
	newPage.CreatedAt = currentPage.CreatedAt
	newPage.UpdatedAt = time.Now()

	return newPage
}

func (pv *PageValidator) markSimilarPagesAsUnnotified(page *model.Page) error {
	similarPages, err := pv.repository.GetPageAndSimilarPages(page)
	if err != nil {
		pv.logger.Error("Error fetching similar pages", "pageID", page.ID, "err", err)
		return err
	}

	for _, similar := range similarPages {
		similar.Notified = false
	}

	return pv.repository.UpdateInBatches(similarPages, 100)
}

func (pv *PageValidator) deletePage(page *model.Page) error {
	if err := pv.repository.Delete(page.ID); err != nil {
		pv.logger.Error("Error deleting page", "pageID", page.ID, "err", err)
		return err
	}
	pv.logger.Info("Deleted expired page", "pageID", page.ID)
	return types.ErrRecordNotFound
}
