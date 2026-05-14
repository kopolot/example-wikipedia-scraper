package service

type FilterQueryInterface interface {
	AddCondition(condition string, arg ...any)
	GetQuery() (string, []any)
	CombineQuerys(other FilterQueryInterface)
}

type BaseFilterQuery struct {
	query string
	args  []any
}

func (f *BaseFilterQuery) AddCondition(condition string, arg ...any) {
	if f.query != "" {
		f.query += " AND "
	}
	f.query += condition
	f.args = append(f.args, arg...)
}

func (f *BaseFilterQuery) GetQuery() (string, []any) {
	return f.query, f.args
}

func (f *BaseFilterQuery) CombineQuerys(other FilterQueryInterface) {
	query, args := other.GetQuery()
	if f.query == "" {
		f.query = query
		f.args = args
		return
	}
	if query == "" {
		return
	}
	f.query = "(" + f.query + ") OR (" + query + ")"
	f.args = append(f.args, args...)
}
