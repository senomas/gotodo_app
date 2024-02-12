package service

import (
	"context"
	"database/sql"
	"encoding/json"
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
)

type NullString struct {
	sql.NullString
}

func NewNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: true}}
}

func NewNullStringNil() NullString {
	return NullString{sql.NullString{Valid: false}}
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

func ServiceContext(ctx context.Context) context.Context {
	if todoService == nil {
		panic("todo service not registered")
	}
	return context.WithValue(ctx, TodoServiceContext, todoService)
}
