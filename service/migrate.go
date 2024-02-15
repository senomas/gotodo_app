package service

import (
	"bufio"
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Migration struct {
	Timestamp time.Time
	Filename  string
	Hash      string
	Result    string
	ID        int64
	Success   bool
}

func FileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filename, err)
	}
	defer f.Close()
	hasher := sha512.New()
	p := make([]byte, 1024)
	for {
		n, err := f.Read(p)
		hasher.Write(p[:n])
		if err == io.EOF {
			return fmt.Sprintf("%x", hasher.Sum(nil)), nil
		} else if err != nil {
			return "", err
		}
	}
}

func Migrate(
	ctx context.Context, path string, addMigrate func(context.Context, Migration) error,
	migrateQuery func(context.Context, string) error,
) error {
	slog.Debug("Migrate", "path", path)
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		fp := filepath.Join(path, f.Name())
		hash, err := FileHash(fp)
		if err != nil {
			return err
		}
		m := Migration{
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
				err := migrateQuery(ctx, qry)
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
			err := migrateQuery(ctx, qry)
			if err != nil {
				return fmt.Errorf("error migrating %s: [%s]\n%v", fp, qry, err)
			}
		}
		m.Success = true
		m.Timestamp = time.Now()
		err = addMigrate(ctx, m)
		if err != nil {
			return err
		}
	}
	return nil
}
