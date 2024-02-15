package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	service "github.com/senomas/gotodo_service"
)

// Migrate implements service.TodoService.
func (TodoService) Migrate(ctx context.Context) error {
	if db, ok := ctx.Value(service.ServiceContextDB).(*sql.DB); ok {
		if path, ok := os.LookupEnv("MIGRATION_PATH"); ok {
			qry := `
        CREATE TABLE IF NOT EXISTS _migration (
          id        INTEGER PRIMARY KEY AUTOINCREMENT,
          filename  TEXT,
          hash      TEXT,
          success   BOOLEAN,
          result    TEXT,
          timestamp DATETIME
        )
      `
			_, err := db.ExecContext(ctx, qry)
			if err != nil {
				slog.Warn("sql error", "qry", qry, "error", err)
				return err
			}
			if path == "" {
				_, filename, _, _ := runtime.Caller(0)
				path = filepath.Join(filepath.Dir(filename), "migration")
			} else if !strings.HasPrefix(path, "/") {
				ex, err := os.Executable()
				if err != nil {
					return err
				}
				path = filepath.Join(ex, path)
			} else {
				path = filepath.Clean(path)
			}
			err = service.Migrate(ctx, path, func(ctx context.Context, m service.Migration) error {
				qry = `
          INSERT INTO _migration (filename, hash, success, result, timestamp)
          VALUES ($1, $2, $3, $4, $5)
        `
				rs, err := db.ExecContext(ctx, qry, m.Filename, m.Hash, m.Success, m.Result, m.Timestamp)
				if err != nil {
					return err
				}
				_, err = rs.LastInsertId()
				if err != nil {
					return err
				}
				return nil
			}, func(ctx context.Context, qry string) error {
				_, err := db.ExecContext(ctx, qry)
				return err
			})
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(`
        PRAGMA foreign_keys = ON;
        PRAGMA integrity_check;
      `)
			if err != nil {
				return err
			}

			_, err = db.Exec(`
      CREATE TABLE IF NOT EXISTS todo_category (
          id INTEGER PRIMARY KEY,
          name TEXT NOT NULL
        )
      `)
			if err != nil {
				return err
			}

			_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS todo (
          id INTEGER PRIMARY KEY,
          title TEXT NOT NULL,
          description TEXT,
          category_id INTEGER NOT NULL,
          done BOOLEAN NOT NULL DEFAULT FALSE,
          FOREIGN KEY (category_id) REFERENCES todo_category (id)
        )
      `)
			return err
		}
	} else {
		return service.ErrNoDBInContext
	}
	return nil
}
