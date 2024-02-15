package sqlite

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	service "github.com/senomas/gotodo_service"
)

// Migrate implements service.TodoService.
func (ts *TodoService) Migrate(ctx context.Context) error {
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
				slog.Warn("Error insert todo", "qry", qry, "error", err)
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
			slog.Debug("Migrate", "path", path)
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, f := range files {
				fp := filepath.Join(path, f.Name())
				hash, err := service.FileHash(fp)
				if err != nil {
					return err
				}
				m := &service.Migration{
					Filename: f.Name(),
					Hash:     hash,
					Result:   "",
					Success:  false,
				}

				fin, err := os.Open(fp)
				if err != nil {
					return fmt.Errorf("error reading %s: %v", fp, err)
				}
				defer fin.Close()
				slog.Debug("Migrate", "file", fp)
				scanner := bufio.NewScanner(fin)
				qry := ""
				for scanner.Scan() {
					ln := scanner.Text()
					qry = fmt.Sprintf("%s%s\n", qry, ln)
					if strings.HasSuffix(strings.TrimSpace(ln), ";") {
						err := ts.migrateQuery(ctx, db, qry)
						if err != nil {
							m.Result = fmt.Sprintf("%s%s\nERROR: %v\n", m.Result, qry, err)
							return fmt.Errorf("error migrating %s: [%s]\n%v", fp, qry, err)
						} else {
							m.Result = fmt.Sprintf("%s%s\n", m.Result, qry)
						}
						qry = ""
					}
				}
				if strings.TrimSpace(qry) != "" {
					err := ts.migrateQuery(ctx, db, qry)
					if err != nil {
						return fmt.Errorf("error migrating %s: [%s]\n%v", fp, qry, err)
					}
				}
				m.Success = true
				m.Timestamp = time.Now()
				qry = `
          INSERT INTO _migration (filename, hash, success, result, timestamp)
          VALUES ($1, $2, $3, $4, $5)
        `
				rs, err := db.ExecContext(ctx, qry, m.Filename, m.Hash, m.Success, m.Result, m.Timestamp)
				if err != nil {
					return err
				}
				m.ID, err = rs.LastInsertId()
				if err != nil {
					return err
				}
				slog.Debug("migrate", "v", m)
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

func (TodoService) migrateQuery(ctx context.Context, db *sql.DB, qry string) error {
	_, err := db.ExecContext(ctx, qry)
	return err
}
