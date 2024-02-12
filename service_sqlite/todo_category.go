package sqlite

import (
	"context"
	"database/sql"

	service "github.com/senomas/gotodo_service"
)

// CreateCategory implements service.TodoService.
func (TodoService) CreateCategory(ctx context.Context, categories []service.TodoCategory) ([]int64, error) {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
		stmt, err := tx.Prepare("INSERT INTO todo_category (name) VALUES (?)")
		if err != nil {
			return nil, err
		}
		ids := make([]int64, len(categories))
		for i, category := range categories {
			res, err := stmt.Exec(category.Name)
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

// DeleteCategory implements service.TodoService.
func (TodoService) DeleteCategory(ctx context.Context, ids []int64) error {
	panic("unimplemented")
}

// UpdateCategory implements service.TodoService.
func (TodoService) UpdateCategory(ctx context.Context, categories []service.TodoCategory) error {
	panic("unimplemented")
}
