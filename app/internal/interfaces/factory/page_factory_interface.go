package factory

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
)

type PageFactoryInterface interface {
	CreateFromDTO(pageDTO *dto.PageDTO) (*model.Page, error)
}
