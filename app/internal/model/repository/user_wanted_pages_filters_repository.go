package repository

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/pkg/repository"
)

type UserWantedPagesFilterRepository struct {
	*repository.GenericRepository[*model.UserWantedPagesFilter]
}

func NewUserWantedPagesFilterRepository() *UserWantedPagesFilterRepository {
	return &UserWantedPagesFilterRepository{
		GenericRepository: repository.NewGenericRepository[*model.UserWantedPagesFilter](),
	}
}
