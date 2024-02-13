package sqlite

import service "github.com/senomas/gotodo_service"

type FilterBool struct {
	field string
	query QueryBuilder
}

// generate implements service.Filter.
func (f *FilterBool) Generate(query service.QueryBuilder) {
	query.AddQuery(&f.query)
}

// Equal implements service.FilterBool.
func (*FilterBool) Equal(bool) service.Filter {
	panic("unimplemented")
}
