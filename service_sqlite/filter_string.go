package sqlite

import service "github.com/senomas/gotodo_service"

type FilterString struct {
	field string
	query QueryBuilder
}

// generate implements service.Filter.
func (f *FilterString) Generate(query service.QueryBuilder) {
	query.AddQuery(&f.query)
}

// Equal implements service.FilterString.
func (f *FilterString) Equal(v string) service.Filter {
	f.query.AddTextParams(f.field+" = ?", v)
	return f
}

// Like implements service.FilterString.
func (f *FilterString) Like(v string) service.Filter {
	f.query.AddTextParams(f.field+" like ?", v)
	return f
}

// NotEqual implements service.FilterString.
func (f *FilterString) NotEqual(string) service.Filter {
	panic("unimplemented")
}

// NotLike implements service.FilterString.
func (f *FilterString) NotLike(string) service.Filter {
	panic("unimplemented")
}

// In implements service.FilterString.
func (f *FilterString) In(v []string) service.Filter {
	str := f.field + " in ("
	params := make([]any, len(v))
	for i, s := range v {
		if i > 0 {
			str += ","
		}
		str += "?"
		params[i] = s
	}
	str += ")"
	f.query.AddTextParams(str, params...)
	return f
}
