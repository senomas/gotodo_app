package sqlite

import (
	"strings"

	service "github.com/senomas/gotodo_service"
)

type QueryBuilder struct {
	prefix string
	sep    string
	sql    []string
	params []any
}

// Sql implements service.QueryBuilder.
func (*QueryBuilder) Sql() string {
	panic("unimplemented")
}

// SqlWithSeparator implements service.QueryBuilder.
func (*QueryBuilder) SqlWithSeparator(sep string) string {
	panic("unimplemented")
}

// Params implements service.QueryBuilder.
func (q *QueryBuilder) Params() []any {
	return q.params
}

// AddParam implements service.QueryBuilder.
func (q *QueryBuilder) AddParam(param any) {
	q.params = append(q.params, param)
}

// AddParams implements service.QueryBuilder.
func (q *QueryBuilder) AddParams(params ...any) {
	q.params = append(q.params, params...)
}

// AddQuery implements service.QueryBuilder.
func (q *QueryBuilder) AddQuery(query service.QueryBuilder) {
	sql := query.SQL()
	if sql != "" {
		q.sql = append(q.sql, sql)
	}
	q.params = append(q.params, query.Params()...)
}

// AddText implements service.QueryBuilder.
func (q *QueryBuilder) AddText(text string) {
	q.sql = append(q.sql, text)
}

// AddTextParam implements service.QueryBuilder.
func (q *QueryBuilder) AddTextParam(text string, param any) {
	q.sql = append(q.sql, text)
	q.params = append(q.params, param)
}

// AddTextParams implements service.QueryBuilder.
func (q *QueryBuilder) AddTextParams(text string, params ...any) {
	q.sql = append(q.sql, text)
	q.params = append(q.params, params...)
}

// SQL implements service.QueryBuilder.
func (q *QueryBuilder) SQL() string {
	if len(q.sql) > 0 {
		return q.prefix + strings.Join(q.sql, q.sep)
	}
	return ""
}
