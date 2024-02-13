package service

import (
	"errors"
)

type ServiceContextType int

const (
	ServiceContextDB ServiceContextType = iota
	ServiceContextCache
	FilterServiceContext
	TodoServiceContext
)

var (
	ErrNoDBInContext = errors.New("DB not found in context")
	ErrNoData        = errors.New("no data")
	ErrInvalidFilter = errors.New("invalid filter")
)
