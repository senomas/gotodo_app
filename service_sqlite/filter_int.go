package sqlite

import service "github.com/senomas/gotodo_service"

type FilterInt struct {
	field string
	query QueryBuilder
}

// generate implements service.Filter.
func (f *FilterInt) Generate(query service.QueryBuilder) {
	query.AddQuery(&f.query)
}

// Between implements service.FilterInt.
func (f *FilterInt) Between(v1 int64, v2 int64) service.Filter {
	f.query.AddTextParams(f.field+" > ? AND "+f.field+" < ?", v1, v2)
	return f
}

// Equal implements service.FilterInt.
func (*FilterInt) Equal(int64) service.Filter {
	panic("unimplemented")
}

// Greater implements service.FilterInt.
func (*FilterInt) Greater(int64) service.Filter {
	panic("unimplemented")
}

// GreaterOrEqual implements service.FilterInt.
func (*FilterInt) GreaterOrEqual(int64) service.Filter {
	panic("unimplemented")
}

// Less implements service.FilterInt.
func (*FilterInt) Less(int64) service.Filter {
	panic("unimplemented")
}

// LessOrEqual implements service.FilterInt.
func (*FilterInt) LessOrEqual(int64) service.Filter {
	panic("unimplemented")
}

// NotEqual implements service.FilterInt.
func (*FilterInt) NotEqual(int64) service.Filter {
	panic("unimplemented")
}
