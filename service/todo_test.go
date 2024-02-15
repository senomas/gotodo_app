package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	service "github.com/senomas/gotodo_service"
	service_impl "github.com/senomas/gotodo_service_sqlite"
	"github.com/stretchr/testify/assert"
)

func init() {
	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(log)
	os.Setenv("MIGRATION_PATH", "")
}

func TestCrud(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	assert.NoError(t, err, "failed to open db")
	defer db.Close()

	ctx := service_impl.NewContext(context.WithValue(context.Background(), service.ServiceContextDB, db))
	todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
	todoService.Migrate(ctx)

	t.Run("CreateCategory", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		categories := []service.TodoCategory{
			{Name: "category 1"},
			{Name: "category 2"},
		}
		ids, err := todoService.CreateCategory(ctx, categories)
		assert.NoError(t, err)
		assert.EqualValues(t, []int64{1, 2}, ids)
	})

	t.Run("Create", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		todos := []service.Todo{
			{
				Category: service.TodoCategory{ID: 1},
				Title:    "todo 1",
			},
			{
				Category:    service.TodoCategory{ID: 1},
				Title:       "todo 2",
				Description: sql.NullString{String: "desc 2", Valid: true},
			},
			{
				Category: service.TodoCategory{ID: 2},
				Title:    "todo 3",
			},
		}
		ids, err := todoService.Create(ctx, todos)
		assert.NoError(t, err)
		assert.EqualValues(t, []int64{1, 2, 3}, ids)
	})

	t.Run("Get", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		todo, err := todoService.Get(ctx, 1)
		assert.NoError(t, err)
		assert.EqualValues(t, service.Todo{
			ID:    1,
			Title: "todo 1",
			Category: service.TodoCategory{
				ID:   1,
				Name: "category 1",
			},
		}, todo)
	})

	getTodoID := func(v any) any {
		if todo, ok := v.(service.Todo); ok {
			return int(todo.ID)
		}
		return "not todo"
	}
	getTodoTitle := func(v any) any {
		if todo, ok := v.(service.Todo); ok {
			return todo.Title
		}
		return "not todo"
	}
	getTodoDescription := func(v any) any {
		if todo, ok := v.(service.Todo); ok {
			if todo.Description.Valid {
				return todo.Description.String
			}
			return nil
		}
		return "not todo"
	}
	getTodoCategory := func(v any) any {
		if todo, ok := v.(service.Todo); ok {
			return todo.Category.Name
		}
		return "not todo"
	}

	t.Run("Find", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		total, todos, err := todoService.Find(ctx, nil, 0, 10)
		assert.NoError(t, err)
		assert.EqualValues(t, 3, total)
		assert.EqualValues(t,
			[]any{1, 2, 3},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 1", "todo 2", "todo 3"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{nil, "desc 2", nil},
			Apply(todos, getTodoDescription))
		assert.EqualValues(t,
			[]any{"category 1", "category 1", "category 2"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Update", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		todos := []service.Todo{
			{
				ID:          3,
				Category:    service.TodoCategory{ID: 1},
				Title:       "todo tiga",
				Description: sql.NullString{String: "desc 3", Valid: true},
			},
		}
		count, err := todoService.Update(ctx, todos)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, count)
	})

	t.Run("Find updated", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		total, todos, err := todoService.Find(ctx, nil, 0, 10)
		assert.NoError(t, err)
		assert.EqualValues(t, 3, total)
		assert.EqualValues(t,
			[]any{1, 2, 3},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 1", "todo 2", "todo tiga"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{nil, "desc 2", "desc 3"},
			Apply(todos, getTodoDescription))
		assert.EqualValues(t,
			[]any{"category 1", "category 1", "category 1"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Create 113 records", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		todos := []service.Todo{}
		eids := []int64{}
		for i := 4; i <= 113; i++ {
			todos = append(todos, service.Todo{
				Title:    fmt.Sprintf("todo %d", i),
				Category: service.TodoCategory{ID: int64(((i / 3) % 2) + 1)},
			})
			eids = append(eids, int64(i))
		}

		ids, err := todoService.Create(ctx, todos)
		assert.NoError(t, err)
		assert.EqualValues(t, eids, ids)
	})

	t.Run("Find with offset limit", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		total, todos, err := todoService.Find(ctx, nil, 4, 5)
		assert.NoError(t, err)
		assert.EqualValues(t, 113, total)
		assert.EqualValues(t,
			[]any{5, 6, 7, 8, 9},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 5", "todo 6", "todo 7", "todo 8", "todo 9"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{"category 2", "category 1", "category 1", "category 1", "category 2"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Find with filter title like", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Title().Like("%11%")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 5, total)
		assert.EqualValues(t,
			[]any{11, 110},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 11", "todo 110"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{"category 2", "category 1"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Find with filter category.name eq", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Category().Equal("category 2")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 56, total)
		assert.EqualValues(t,
			[]any{4, 5},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 4", "todo 5"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{"category 2", "category 2"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Find with multiple filter", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Title().Like("%11%")
		filter.Category().Equal("category 2")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 4, total)
		assert.EqualValues(t,
			[]any{11, 111},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 11", "todo 111"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{"category 2", "category 2"},
			Apply(todos, getTodoCategory))
	})

	t.Run("Find with filter safe in", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Category().In([]string{"category 2", `"safe" in`})
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 56, total)
		assert.EqualValues(t,
			[]any{4, 5},
			Apply(todos, getTodoID))
		assert.EqualValues(t,
			[]any{"todo 4", "todo 5"},
			Apply(todos, getTodoTitle))
		assert.EqualValues(t,
			[]any{"category 2", "category 2"},
			Apply(todos, getTodoCategory))
	})
}

func Apply(value any, fn func(v any) any) []any {
	va := reflect.ValueOf(value)
	res := make([]any, va.Len())
	for i := 0; i < va.Len(); i++ {
		res[i] = fn(va.Index(i).Interface())
	}
	return res
}
