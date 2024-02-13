package sqlite_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	service "github.com/senomas/gotodo_service"
	service_impl "github.com/senomas/gotodo_service_sqlite"
	"github.com/stretchr/testify/assert"
)

func TestCrud(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	assert.NoError(t, err, "failed to open db")
	defer db.Close()

	_, err = db.Exec(`
    PRAGMA foreign_keys = ON;
    PRAGMA integrity_check;
  `)
	assert.NoError(t, err, "failed to enable foreign keys")

	_, err = db.Exec(`
   CREATE TABLE todo_category (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL
    )
  `)
	assert.NoError(t, err, "failed to create table")

	_, err = db.Exec(`
    CREATE TABLE todo (
      id INTEGER PRIMARY KEY,
      title TEXT NOT NULL,
      description TEXT,
      category_id INTEGER NOT NULL,
      done BOOLEAN NOT NULL DEFAULT FALSE,
      FOREIGN KEY (category_id) REFERENCES todo_category (id)
    )
  `)
	assert.NoError(t, err, "failed to create table")

	ctx := service_impl.NewContext(context.WithValue(context.Background(), service.ServiceContextDB, db))

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
				Description: service.NewNullString("desc 2"),
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

	t.Run("Find", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		total, todos, err := todoService.Find(ctx, nil, 0, 10)
		assert.NoError(t, err)
		assert.EqualValues(t, 3, total)
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 1",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 1,
        "done": false
      },
      {
        "title": "todo 2",
        "description": "desc 2",
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 2,
        "done": false
      },
      {
        "title": "todo 3",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 3,
        "done": false
      }
    ]`, string(str))
	})

	t.Run("Update", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		todos := []service.Todo{
			{
				ID:          3,
				Category:    service.TodoCategory{ID: 1},
				Title:       "todo tiga",
				Description: service.NewNullString("desc 3"),
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
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 1",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 1,
        "done": false
      },
      {
        "title": "todo 2",
        "description": "desc 2",
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 2,
        "done": false
      },
      {
        "title": "todo tiga",
        "description": "desc 3",
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 3,
        "done": false
      }
    ]`, string(str))
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
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 5",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 5,
        "done": false
      },
      {
        "title": "todo 6",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 6,
        "done": false
      },
      {
        "title": "todo 7",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 7,
        "done": false
      },
      {
        "title": "todo 8",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 8,
        "done": false
      },
      {
        "title": "todo 9",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 9,
        "done": false
      }
    ]`, string(str))
	})

	t.Run("Find with filter title like", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Title().Like("%11%")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 5, total)
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 11",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 11,
        "done": false
      },
      {
        "title": "todo 110",
        "description": null,
        "category": {
          "name": "category 1",
          "id": 1
        },
        "id": 110,
        "done": false
      }
    ]`, string(str))
	})

	t.Run("Find with filter category.name eq", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Category().Equal("category 2")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 56, total)
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 4",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 4,
        "done": false
      },
      {
        "title": "todo 5",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 5,
        "done": false
      }
    ]`, string(str))
	})

	t.Run("Find with multiple filter", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Title().Like("%11%")
		filter.Category().Equal("category 2")
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 4, total)
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 11",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 11,
        "done": false
      },
      {
        "title": "todo 111",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 111,
        "done": false
      }
    ]`, string(str))
	})

	t.Run("Find with filter safe in", func(t *testing.T) {
		todoService := ctx.Value(service.TodoServiceContext).(service.TodoService)
		filter := todoService.Filter()
		filter.Category().In([]string{"category 2", `"safe" in`})
		total, todos, err := todoService.Find(ctx, filter, 0, 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 56, total)
		str, err := json.MarshalIndent(todos, "    ", "  ")
		assert.NoError(t, err)
		assert.Equal(t, `[
      {
        "title": "todo 4",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 4,
        "done": false
      },
      {
        "title": "todo 5",
        "description": null,
        "category": {
          "name": "category 2",
          "id": 2
        },
        "id": 5,
        "done": false
      }
    ]`, string(str))
	})
}
