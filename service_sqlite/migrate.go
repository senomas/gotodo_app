package sqlite

import (
	"context"
	"database/sql"

	service "github.com/senomas/gotodo_service"
)

// Migrate implements service.TodoService.
func (TodoService) Migrate(ctx context.Context) error {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		_, err := db.Exec(`
      PRAGMA foreign_keys = ON;
      PRAGMA integrity_check;
    `)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
    CREATE TABLE todo_category (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL
      )
    `)
		if err != nil {
			return err
		}

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
		return err
	} else {
		return service.ErrNoDBInContext
	}
}
