package sqlite

import (
	"context"

	service "github.com/senomas/gotodo_service"
)

func NewContext(ctx context.Context) context.Context {
	var todoService service.TodoService = &TodoService{}
	return context.WithValue(
		ctx,
		service.TodoServiceContext,
		todoService,
	)
}
