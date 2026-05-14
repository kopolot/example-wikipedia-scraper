package repository

import (
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/pkg/db"
)

type RepositoryInterface[T interfaces.ModelInterface] interface {
	GetAll() ([]T, error)
	GetAllWithPreloads(preloads ...string) ([]T, error)
	GetByID(id uint) (T, error)
	GetByIDWithPreloads(id uint, preloads ...string) (T, error)
	Create(model T) error
	Upsert(model T, columns ...string) error
	Update(model T) error
	Delete(id uint) error
	Disable(id uint) error
	CreateInBatches(models []T, batchSize uint) error
	UpsertInBatches(models []T, batchSize uint, columns ...string) error
	UpdateInBatches(models []T, batchSize uint) error
	FindBy(conds ...any) ([]T, error)
	FindOneBy(conds ...any) (T, error)
	FindByWithLimit(conds []any, limit uint) ([]T, error)
	FindByWithLimitPage(conds []any, limit uint, page uint) ([]T, error)
	FindByWithLimitPageWithPreloads(conds []any, limit uint, page uint) ([]T, error)
	CountBy(query any, args ...any) (int64, error)
	GetQueryBuilder() db.QueryBuilder
}
