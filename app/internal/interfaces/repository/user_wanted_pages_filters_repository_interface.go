package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
)

type UserWantedPagesFilterRepositoryInterface interface {
	repository.RepositoryInterface[*model.UserWantedPagesFilter]
}
