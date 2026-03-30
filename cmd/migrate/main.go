package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"WITS/config"
	"WITS/package/database"

	"github.com/joho/godotenv"
)

type migrationFile struct {
	Name string
	Path string
}

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := collectMigrationFiles("migration")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		sql, err := os.ReadFile(file.Path)
		if err != nil {
			log.Fatalf("read %s: %v", file.Name, err)
		}

		if _, err := db.Exec(context.Background(), string(sql)); err != nil {
			log.Fatalf("apply %s: %v", file.Name, err)
		}

		log.Printf("applied %s", file.Name)
	}

	log.Println("migrations completed successfully")
}

func collectMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	files := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}

		files = append(files, migrationFile{
			Name: entry.Name(),
			Path: filepath.Join(dir, entry.Name()),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return migrationPriority(files[i].Name) < migrationPriority(files[j].Name)
	})

	return files, nil
}

func migrationPriority(name string) int {
	switch name {
	case "005_asset_seq.up.sql":
		return 1
	case "005_asset_inventory.up.sql":
		return 2
	case "005_asset_assignments.up.sql":
		return 3
	case "005_asset_maintenance.up.sql":
		return 4
	default:
		return 100
	}
}
