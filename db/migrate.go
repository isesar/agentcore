package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(migrationsPath string) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := ensureMigrationsTable(DB); err != nil {
		return err
	}

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	upMigrations := filterAndSortUpMigrations(files)
	for _, migration := range upMigrations {
		applied, err := isMigrationApplied(DB, migration.Name())
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsPath, migration.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration.Name(), err)
		}

		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", migration.Name(), err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Name(), err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, migration.Name()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to register migration %s: %w", migration.Name(), err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Name(), err)
		}
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}
	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`,
		version,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check migration %s: %w", version, err)
	}
	return exists, nil
}

func filterAndSortUpMigrations(files []fs.DirEntry) []fs.DirEntry {
	migrations := make([]fs.DirEntry, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrations = append(migrations, file)
		}
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})
	return migrations
}
