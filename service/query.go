package service

type QueryBuilder interface {
	AddText(text string)
	AddTextParam(text string, param any)
	AddTextParams(text string, params ...any)
	AddParam(param any)
	AddParams(params ...any)

	AddQuery(query QueryBuilder)

	SQL() string
	Params() []any
}
