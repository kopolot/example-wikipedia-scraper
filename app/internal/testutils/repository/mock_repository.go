package testutils

import (
	"example-wikipedia-scraper/pkg/db"

	"github.com/stretchr/testify/mock"
)

type MockGenericRepository[T any] struct {
	mock.Mock
}

func (m *MockGenericRepository[T]) GetAll() ([]T, error) {
	args := m.Called()
	return args.Get(0).([]T), args.Error(1)
}
func (m *MockGenericRepository[T]) GetAllWithPreloads(preloads ...string) ([]T, error) {
	args := m.Called(preloads)
	return args.Get(0).([]T), args.Error(1)
}
func (m *MockGenericRepository[T]) GetByID(id uint) (T, error) {
	args := m.Called(id)
	return args.Get(0).(T), args.Error(1)
}
func (m *MockGenericRepository[T]) GetByIDWithPreloads(id uint, preloads ...string) (T, error) {
	args := m.Called(id, preloads)
	return args.Get(0).(T), args.Error(1)
}
func (m *MockGenericRepository[T]) Create(model T) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) Upsert(model T, columns ...string) error {
	args := m.Called(model, columns)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) Update(model T) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) Disable(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) CreateInBatches(models []T, batchSize uint) error {
	args := m.Called(models, batchSize)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) UpsertInBatches(models []T, batchSize uint, columns ...string) error {
	args := m.Called(models, batchSize, columns)
	return args.Error(0)
}
func (m *MockGenericRepository[T]) FindBy(conds ...any) ([]T, error) {
	args := m.Called(conds)
	return args.Get(0).([]T), args.Error(1)
}
func (m *MockGenericRepository[T]) FindOneBy(conds ...any) (T, error) {
	args := m.Called(conds)
	return args.Get(0).(T), args.Error(1)
}
func (m *MockGenericRepository[T]) FindByWithLimit(conds []any, limit uint) ([]T, error) {
	args := m.Called(conds, limit)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockGenericRepository[T]) FindByWithLimitPage(conds []any, limit uint, page uint) ([]T, error) {
	args := m.Called(conds, limit, page)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockGenericRepository[T]) FindByWithLimitPageWithPreloads(conds []any, limit uint, page uint) ([]T, error) {
	args := m.Called(conds, limit, page)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockGenericRepository[T]) CountBy(query any, args ...any) (int64, error) {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(int64), mockArgs.Error(1)
}

func (m *MockGenericRepository[T]) UpdateInBatches(models []T, batchSize uint) error {
	args := m.Called(models, batchSize)
	return args.Error(0)
}

func (m *MockGenericRepository[T]) GetQueryBuilder() db.QueryBuilder {
	args := m.Called()
	return args.Get(0).(db.QueryBuilder)
}
