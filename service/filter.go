package service

type FilterOp int

type FilterString interface {
	Any()
	Equal(string)
	NotEqual(string)
	Like(string)
	NotLike(string)
}

type FilterInt interface {
	Any()
	Equal(int64)
	NotEqual(int64)
	Less(int64)
	LessOrEqual(int64)
	Greater(int64)
	GreaterOrEqual(int64)
	Between(int64, int64)
}

type FilterBool interface {
	Any()
	Equal(bool)
}

type FilterService interface {
	FilterString() FilterString
	FilterInt() FilterInt
	FilterBool() FilterBool
}

var filterService FilterService

func RegisterFilterService(s FilterService) {
	if filterService != nil {
		panic("filter service already registered")
	}
	filterService = s
}
