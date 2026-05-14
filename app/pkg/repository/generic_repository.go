package repository

import (
	"example-wikipedia-scraper/internal/db"
	"example-wikipedia-scraper/internal/interfaces"
	pkgDb "example-wikipedia-scraper/pkg/db"
	"slices"
)

type GenericRepository[T interfaces.ModelInterface] struct {
	QB pkgDb.QueryBuilder
}

func NewGenericRepository[T interfaces.ModelInterface]() *GenericRepository[T] {
	return &GenericRepository[T]{QB: pkgDb.NewQueryBuilder(db.DB)}
}

func (r *GenericRepository[T]) GetAll() ([]T, error) {
	var models []T
	err := r.QB.Find(&models)
	return handleNotFoundSlice(err, models)
}

func (r *GenericRepository[T]) GetAllWithPreloads(preloads ...string) ([]T, error) {
	var models []T
	qb := r.QB
	if len(preloads) == 0 {
		qb = qb.PreloadAssociations()
	} else {
		for _, preload := range preloads {
			qb = qb.Preload(preload)
		}
	}
	err := qb.Find(&models)
	return handleNotFoundSlice(err, models)
}

func (r *GenericRepository[T]) GetByID(id uint) (T, error) {
	var model T
	err := r.QB.Unscoped().First(&model, id)
	return handleNotFoundModel(err, model)
}

func (r *GenericRepository[T]) GetByIDWithPreloads(id uint, preloads ...string) (T, error) {
	var model T
	qb := r.QB
	if len(preloads) == 0 {
		qb = qb.PreloadAssociations()
	} else {
		for _, preload := range preloads {
			qb = qb.Preload(preload)
		}
	}
	err := qb.Unscoped().First(&model, id)
	return handleNotFoundModel(err, model)
}

func (r *GenericRepository[T]) Create(model T) error {
	return r.QB.Create(model)
}

func (r *GenericRepository[T]) Upsert(model T, columns ...string) error {
	return r.QB.Upsert(model, columns...)
}

func (r *GenericRepository[T]) Update(model T) error {
	return r.QB.Save(model)
}

func (r *GenericRepository[T]) Delete(id uint) error {
	var model T
	return r.QB.Unscoped().Delete(&model, id)
}

func (r *GenericRepository[T]) Disable(id uint) error {
	var model T
	return r.QB.Delete(&model, id)
}

func (r *GenericRepository[T]) CreateInBatches(models []T, batchSize uint) error {
	if len(models) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	tx := r.QB.Begin()
	iter := slices.Chunk(models, int(batchSize))
	for batch := range iter {
		if err := tx.Create(&batch); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *GenericRepository[T]) UpsertInBatches(models []T, batchSize uint, columns ...string) error {
	if len(models) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	tx := r.QB.Begin()
	iter := slices.Chunk(models, int(batchSize))
	for batch := range iter {
		if err := tx.Upsert(&batch, columns...); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *GenericRepository[T]) UpdateInBatches(models []T, batchSize uint) error {
	if len(models) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	tx := r.QB.Begin()
	iter := slices.Chunk(models, int(batchSize))
	for batch := range iter {
		if err := tx.Save(&batch); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *GenericRepository[T]) FindBy(conds ...any) ([]T, error) {
	var models []T
	err := r.QB.Find(&models, conds...)
	return handleNotFoundSlice(err, models)
}

func (r *GenericRepository[T]) FindOneBy(conds ...any) (T, error) {
	var model T
	err := r.QB.First(&model, conds...)
	return handleNotFoundModel(err, model)
}

func (r *GenericRepository[T]) FindByWithLimit(conds []any, limit uint) ([]T, error) {
	return r.FindByWithLimitPage(conds, limit, 1)
}

func (r *GenericRepository[T]) FindByWithLimitPage(conds []any, limit uint, page uint) ([]T, error) {
	var models []T
	offset := (page - 1) * limit
	err := r.QB.Limit(int(limit)).Offset(int(offset)).Find(&models, conds...)
	return handleNotFoundSlice(err, models)
}

func (r *GenericRepository[T]) FindByWithLimitPageWithPreloads(conds []any, limit uint, page uint) ([]T, error) {
	var models []T
	offset := (page - 1) * limit
	err := r.QB.PreloadAssociations().Limit(int(limit)).Offset(int(offset)).Find(&models, conds...)
	return handleNotFoundSlice(err, models)
}

func (r *GenericRepository[T]) CountBy(query any, args ...any) (int64, error) {
	var count int64
	err := r.QB.Model(new(T)).Where(query, args...).Count(&count)
	return count, err
}

func (r *GenericRepository[T]) GetQueryBuilder() pkgDb.QueryBuilder {
	return r.QB
}
