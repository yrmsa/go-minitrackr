package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Enable WAL mode
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}

	// Memory optimizations
	if _, err := db.Exec("PRAGMA cache_size=-2000"); err != nil { // 2MB cache
		return nil, err
	}
	if _, err := db.Exec("PRAGMA mmap_size=0"); err != nil { // Disable mmap
		return nil, err
	}
	if _, err := db.Exec("PRAGMA temp_store=MEMORY"); err != nil {
		return nil, err
	}

	// Connection pool limits
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := createSchema(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS issues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'todo',
		priority TEXT NOT NULL DEFAULT 'medium',
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);

	CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
	CREATE INDEX IF NOT EXISTS idx_issues_created_at ON issues(created_at DESC);
	`

	_, err := db.Exec(schema)
	return err
}
