package service

type Filter interface {
	Generate(query QueryBuilder)
}

type FilterString interface {
	Equal(string) Filter
	NotEqual(string) Filter
	Like(string) Filter
	NotLike(string) Filter
	In([]string) Filter
}

type FilterInt interface {
	Equal(int64) Filter
	NotEqual(int64) Filter
	Less(int64) Filter
	LessOrEqual(int64) Filter
	Greater(int64) Filter
	GreaterOrEqual(int64) Filter
	Between(int64, int64) Filter
}

type FilterBool interface {
	Equal(bool) Filter
}

type FilterService interface {
	FilterString() FilterString
	FilterInt() FilterInt
	FilterBool() FilterBool
}
