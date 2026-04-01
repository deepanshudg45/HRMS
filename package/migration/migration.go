package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type RunnerDB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type File struct {
	Name string
	Path string
}

func Run(ctx context.Context, db RunnerDB, dir string) error {
	files, err := CollectFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		sql, err := os.ReadFile(file.Path)
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name, err)
		}

		if _, err := db.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply %s: %w", file.Name, err)
		}
	}

	return nil
}

func CollectFiles(dir string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	files := make([]File, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		files = append(files, File{
			Name: name,
			Path: filepath.Join(dir, name),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return priority(files[i].Name) < priority(files[j].Name)
	})

	return files, nil
}

func priority(name string) int {
	switch name {
	case "005_asset_inventory.up.sql":
		return 1
	case "005_asset_assignments.up.sql":
		return 2
	case "005_asset_maintenance.up.sql":
		return 3
	default:
		return 100
	}
}
