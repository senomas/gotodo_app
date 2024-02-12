package sqlite

import (
	"log/slog"

	service "github.com/senomas/gotodo_service"
)

func init() {
	slog.Debug("RegisterTodoService", "implements", "sqlite")
	service.RegisterTodoService(TodoService{})
}
