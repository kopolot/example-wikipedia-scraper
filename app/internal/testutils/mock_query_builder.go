package testutils

import (
	"example-wikipedia-scraper/pkg/db"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm/clause"
)

type MockQueryBuilder struct {
	mock.Mock
}

func (g *MockQueryBuilder) Table(name string, args ...interface{}) db.QueryBuilder {
	g.Called(name, args)
	return g
}

func (g *MockQueryBuilder) Select(query interface{}, args ...interface{}) db.QueryBuilder {
	g.Called(query, args)
	return g
}

func (g *MockQueryBuilder) Where(query interface{}, args ...interface{}) db.QueryBuilder {
	g.Called(query, args)
	return g
}

func (g *MockQueryBuilder) OrWhere(query interface{}, args ...interface{}) db.QueryBuilder {
	g.Called(query, args)
	return g
}

func (g *MockQueryBuilder) Order(value interface{}) db.QueryBuilder {
	g.Called(value)
	return g
}

func (g *MockQueryBuilder) Limit(limit int) db.QueryBuilder {
	g.Called(limit)
	return g
}

func (g *MockQueryBuilder) Offset(offset int) db.QueryBuilder {
	g.Called(offset)
	return g
}

func (g *MockQueryBuilder) Preload(query string, args ...interface{}) db.QueryBuilder {
	g.Called(query, args)
	return g
}

func (g *MockQueryBuilder) Save(value interface{}) error {
	r := g.Called(value)
	return r.Error(0)
}

func (g *MockQueryBuilder) Create(value interface{}) error {
	r := g.Called(value)
	return r.Error(0)
}

func (g *MockQueryBuilder) Delete(value interface{}, conds ...interface{}) error {
	r := g.Called(value, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) Unscoped() db.QueryBuilder {
	g.Called()
	return g
}

func (g *MockQueryBuilder) Begin() db.QueryBuilder {
	g.Called()
	return g
}

func (g *MockQueryBuilder) Commit() error {
	r := g.Called()
	return r.Error(0)
}

func (g *MockQueryBuilder) Rollback() error {
	r := g.Called()
	return r.Error(0)
}

func (g *MockQueryBuilder) First(dest interface{}, conds ...interface{}) error {
	r := g.Called(dest, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) FirstOrCreate(dest interface{}, conds ...interface{}) error {
	r := g.Called(dest, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) FirstOrInit(dest interface{}, conds ...interface{}) error {
	r := g.Called(dest, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) Find(dest interface{}, conds ...interface{}) error {
	r := g.Called(dest, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) Scan(dest interface{}) error {
	r := g.Called(dest)
	return r.Error(0)
}

func (g *MockQueryBuilder) Pluck(column string, dest interface{}) error {
	r := g.Called(column, dest)
	return r.Error(0)
}

func (g *MockQueryBuilder) Count(count *int64) error {
	r := g.Called(count)
	return r.Error(0)
}

func (g *MockQueryBuilder) Omit(columns ...string) db.QueryBuilder {
	g.Called(columns)
	return g
}

func (g *MockQueryBuilder) Take(dest interface{}, conds ...interface{}) error {
	r := g.Called(dest, conds)
	return r.Error(0)
}

func (g *MockQueryBuilder) Update(column string, value interface{}) error {
	r := g.Called(column, value)
	return r.Error(0)
}

func (g *MockQueryBuilder) Updates(values interface{}) error {
	r := g.Called(values)
	return r.Error(0)
}

func (g *MockQueryBuilder) Exec(sql string, values ...interface{}) error {
	r := g.Called(sql, values)
	return r.Error(0)
}

func (g *MockQueryBuilder) Error() error {
	r := g.Called()
	return r.Error(0)
}

func (g *MockQueryBuilder) RowsAffected() int64 {
	r := g.Called()
	return r.Get(0).(int64)
}

func (g *MockQueryBuilder) PreloadAssociations() db.QueryBuilder {
	g.Called()
	return g
}

func (g *MockQueryBuilder) Clauses(clauses ...clause.Expression) db.QueryBuilder {
	g.Called(clauses)
	return g
}

func (g *MockQueryBuilder) Upsert(value any, columns ...string) error {
	r := g.Called(value, columns)
	return r.Error(0)
}

func (g *MockQueryBuilder) Model(value any) db.QueryBuilder {
	g.Called(value)
	return g
}

func (g *MockQueryBuilder) Distinct(args ...any) db.QueryBuilder {
	g.Called(args)
	return g
}

func (g *MockQueryBuilder) LockForUpdate() db.QueryBuilder {
	g.Called()
	return g
}

func (g *MockQueryBuilder) Group(column string) db.QueryBuilder {
	g.Called(column)
	return g
}

func (g *MockQueryBuilder) Transaction(fn func(tx db.QueryBuilder) error) error {
	args := g.Called(fn)
	return args.Error(0)
}
