package service

import (
	"context"
)

type Todo struct {
	Title       string       `json:"title"`
	Description NullString   `json:"description"`
	Category    TodoCategory `json:"category"`
	ID          int64        `json:"id"`
	Done        bool         `json:"done"`
}

type TodoCategory struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type TodoFilter struct {
	Title       FilterString
	Description FilterString
	Category    FilterString
	CategoryID  FilterInt
	Done        FilterBool
}

type TodoService interface {
	CreateCategory(ctx context.Context, categories []TodoCategory) ([]int64, error)
	UpdateCategory(ctx context.Context, categories []TodoCategory) error
	DeleteCategory(ctx context.Context, ids []int64) error

	Create(ctx context.Context, todos []Todo) ([]int64, error)
	Update(ctx context.Context, todos []Todo) (int64, error)
	Delete(ctx context.Context, ids []int64) error

	Get(ctx context.Context, id int64) (Todo, error)
	Find(
		ctx context.Context, filter TodoFilter, offset int64, limit int,
	) (int64, []Todo, error)
}

var todoService TodoService

func RegisterTodoService(s TodoService) {
	if todoService != nil {
		panic("todo service already registered")
	}
	todoService = s
}
