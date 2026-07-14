package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryBuilder interface {
	Table(name string, args ...interface{}) QueryBuilder
	Select(query interface{}, args ...interface{}) QueryBuilder
	Where(query interface{}, args ...interface{}) QueryBuilder
	OrWhere(query interface{}, args ...interface{}) QueryBuilder
	Order(value interface{}) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Preload(query string, args ...interface{}) QueryBuilder
	Save(value interface{}) error
	// Upsert inserts or updates a record based on the specified columns.
	Upsert(value any, columns ...string) error
	Create(value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	Unscoped() QueryBuilder
	Clauses(clauses ...clause.Expression) QueryBuilder
	Begin() QueryBuilder
	Commit() error
	Rollback() error
	First(dest interface{}, conds ...interface{}) error
	FirstOrCreate(dest interface{}, conds ...interface{}) error
	FirstOrInit(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Scan(dest interface{}) error
	Pluck(column string, dest interface{}) error
	Count(count *int64) error
	Omit(columns ...string) QueryBuilder
	Take(dest interface{}, conds ...interface{}) error
	Update(column string, value interface{}) error
	Updates(values interface{}) error
	Exec(sql string, values ...interface{}) error
	Error() error
	RowsAffected() int64
	PreloadAssociations() QueryBuilder
	Model(value any) QueryBuilder
	Distinct(args ...any) QueryBuilder
	LockForUpdate() QueryBuilder
	Group(column string) QueryBuilder
	Transaction(fn func(tx QueryBuilder) error) error
}

type gormQueryBuilder struct {
	db           *gorm.DB
	rowsAffected int64
}

func NewQueryBuilder(db *gorm.DB) QueryBuilder {
	return &gormQueryBuilder{db: db}
}

func (g *gormQueryBuilder) Table(name string, args ...interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Table(name, args...)}
}

func (g *gormQueryBuilder) Select(query interface{}, args ...interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Select(query, args...)}
}

func (g *gormQueryBuilder) Where(query interface{}, args ...interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Where(query, args...)}
}

func (g *gormQueryBuilder) OrWhere(query interface{}, args ...interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Or(query, args...)}
}

func (g *gormQueryBuilder) Order(value interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Order(value)}
}

func (g *gormQueryBuilder) Limit(limit int) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Limit(limit)}
}

func (g *gormQueryBuilder) Offset(offset int) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Offset(offset)}
}

func (g *gormQueryBuilder) Preload(query string, args ...interface{}) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Preload(query, args...)}
}

func (g *gormQueryBuilder) Save(value interface{}) error {
	return g.db.Save(value).Error
}

func (g *gormQueryBuilder) Create(value interface{}) error {
	return g.db.Create(value).Error
}

func (g *gormQueryBuilder) Delete(value interface{}, conds ...interface{}) error {
	return g.db.Delete(value, conds...).Error
}

func (g *gormQueryBuilder) Unscoped() QueryBuilder {
	return &gormQueryBuilder{db: g.db.Unscoped()}
}

func (g *gormQueryBuilder) Begin() QueryBuilder {
	return &gormQueryBuilder{db: g.db.Begin()}
}

func (g *gormQueryBuilder) Commit() error {
	return g.db.Commit().Error
}

func (g *gormQueryBuilder) Rollback() error {
	return g.db.Rollback().Error
}

func (g *gormQueryBuilder) Transaction(fn func(tx QueryBuilder) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(&gormQueryBuilder{db: tx})
	})
}

func (g *gormQueryBuilder) First(dest interface{}, conds ...interface{}) error {
	return g.db.First(dest, conds...).Error
}

func (g *gormQueryBuilder) FirstOrCreate(dest interface{}, conds ...interface{}) error {
	return g.db.FirstOrCreate(dest, conds...).Error
}

func (g *gormQueryBuilder) FirstOrInit(dest interface{}, conds ...interface{}) error {
	return g.db.FirstOrInit(dest, conds...).Error
}

func (g *gormQueryBuilder) Find(dest interface{}, conds ...interface{}) error {
	return g.db.Find(dest, conds...).Error
}

func (g *gormQueryBuilder) Scan(dest interface{}) error {
	return g.db.Scan(dest).Error
}

func (g *gormQueryBuilder) Pluck(column string, dest interface{}) error {
	return g.db.Pluck(column, dest).Error
}

func (g *gormQueryBuilder) Count(count *int64) error {
	return g.db.Count(count).Error
}

func (g *gormQueryBuilder) Omit(columns ...string) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Omit(columns...)}
}

func (g *gormQueryBuilder) Take(dest interface{}, conds ...interface{}) error {
	return g.db.Take(dest, conds...).Error
}

func (g *gormQueryBuilder) Update(column string, value interface{}) error {
	return g.db.Update(column, value).Error
}

func (g *gormQueryBuilder) Updates(values interface{}) error {
	return g.db.Updates(values).Error
}

func (g *gormQueryBuilder) Exec(sql string, values ...interface{}) error {
	return g.db.Exec(sql, values...).Error
}

func (g *gormQueryBuilder) Error() error {
	return g.db.Error
}

func (g *gormQueryBuilder) RowsAffected() int64 {
	return g.db.RowsAffected
}

func (g *gormQueryBuilder) PreloadAssociations() QueryBuilder {
	return &gormQueryBuilder{db: g.db.Preload(clause.Associations)}
}

func (g *gormQueryBuilder) Clauses(clauses ...clause.Expression) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Clauses(clauses...)}
}

func (g *gormQueryBuilder) Upsert(value any, columns ...string) error {
	var onConflict clause.OnConflict
	if len(columns) > 0 {
		cols := make([]clause.Column, len(columns))
		for i, c := range columns {
			cols[i] = clause.Column{Name: c}
		}
		onConflict.Columns = cols
	}
	onConflict.UpdateAll = true
	return g.db.Clauses(onConflict).Create(value).Error
}

func (g *gormQueryBuilder) Model(value any) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Model(value)}
}

func (g *gormQueryBuilder) Distinct(args ...any) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Distinct(args...)}
}

func (g *gormQueryBuilder) LockForUpdate() QueryBuilder {
	return &gormQueryBuilder{db: g.db.Clauses(clause.Locking{Strength: "UPDATE"})}
}

func (g *gormQueryBuilder) Group(column string) QueryBuilder {
	return &gormQueryBuilder{db: g.db.Group(column)}
}
