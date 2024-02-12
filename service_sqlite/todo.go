package sqlite

import (
	"context"
	"database/sql"
	"log/slog"

	service "github.com/senomas/gotodo_service"
)

type TodoService struct{}

// Create implements service.TodoService.
func (TodoService) Create(ctx context.Context, todos []service.Todo) ([]int64, error) {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		stmt, err := tx.PrepareContext(ctx, "INSERT INTO todo (title, description, category_id, done) VALUES (?, ?, ?, ?)")
		if err != nil {
			return nil, err
		}
		ids := make([]int64, len(todos))
		for i, todo := range todos {
			res, err := stmt.ExecContext(ctx, todo.Title, todo.Description, todo.Category.ID, todo.Done)
			if err != nil {
				return nil, err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return nil, err
			}
			ids[i] = id
		}

		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		return ids, nil
	} else {
		return nil, service.ErrNoDBInContext
	}
}

// Update implements service.TodoService.
func (TodoService) Update(ctx context.Context, todos []service.Todo) (int64, error) {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		tx, err := db.Begin()
		if err != nil {
			return 0, err
		}
		defer tx.Rollback()

		stmt, err := tx.PrepareContext(ctx, "UPDATE todo SET title = ?, description = ?, category_id = ?, done = ? WHERE id = ?")
		if err != nil {
			return 0, err
		}
		var affected int64
		for _, todo := range todos {
			res, err := stmt.ExecContext(ctx, todo.Title, todo.Description, todo.Category.ID, todo.Done, todo.ID)
			if err != nil {
				return 0, err
			}
			aff, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}
			affected += aff
		}

		err = tx.Commit()
		if err != nil {
			return 0, err
		}
		return affected, nil
	} else {
		return 0, service.ErrNoDBInContext
	}
}

// Delete implements service.TodoService.
func (TodoService) Delete(ctx context.Context, ids []int64) error {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		slog.Debug("TodoService.Delete", "db", db)
		panic("unimplemented")
	} else {
		return service.ErrNoDBInContext
	}
}

// Find implements service.TodoService.
func (TodoService) Find(
	ctx context.Context, filter service.TodoFilter, offset int64, limit int,
) (int64, []service.Todo, error) {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		rows, err := db.QueryContext(ctx, `
      SELECT COUNT(t.id)
      FROM todo t JOIN todo_category c ON t.category_id = c.id 
    `, limit, offset)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		var total int64
		if rows.Next() {
			err = rows.Scan(&total)
			if err != nil {
				return 0, nil, err
			}
		}

		rows, err = db.QueryContext(ctx, `
      SELECT t.id, title, description, category_id, c.name, done 
      FROM todo t JOIN todo_category c ON t.category_id = c.id 
      LIMIT ? OFFSET ?
    `, limit, offset)
		if err != nil {
			return total, nil, err
		}
		defer rows.Close()
		var todos []service.Todo
		for rows.Next() {
			var todo service.Todo
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Category.ID, &todo.Category.Name, &todo.Done)
			if err != nil {
				return total, nil, err
			}
			todos = append(todos, todo)
		}
		return total, todos, nil
	} else {
		return 0, nil, service.ErrNoDBInContext
	}
}

// Get implements service.TodoService.
func (TodoService) Get(ctx context.Context, id int64) (service.Todo, error) {
	var todo service.Todo
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		rows, err := db.QueryContext(ctx, `
      SELECT t.id, title, description, category_id, c.name, done 
      FROM todo t JOIN todo_category c ON t.category_id = c.id 
      WHERE t.id = ?
    `, id)
		if err != nil {
			return todo, err
		}
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Category.ID, &todo.Category.Name, &todo.Done)
			return todo, err
		}
		return todo, service.ErrNoData
	} else {
		return todo, service.ErrNoDBInContext
	}
}
