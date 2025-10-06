package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Up applies all pending migrations
func (m *Migrator) Up() error {
	// Ensure migrations table exists
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Apply pending migrations
	for _, file := range migrationFiles {
		name := strings.TrimSuffix(file, ".sql")
		if _, exists := applied[name]; !exists {
			if err := m.applyMigration(file); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file, err)
			}
			m.logger.Info("Applied migration", zap.String("file", file))
		}
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	// Get the last applied migration
	lastMigration, err := m.getLastMigration()
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	if lastMigration == "" {
		m.logger.Info("No migrations to rollback")
		return nil
	}

	// Execute rollback
	if err := m.rollbackMigration(lastMigration); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration, err)
	}

	m.logger.Info("Rolled back migration", zap.String("migration", lastMigration))
	return nil
}

// Create creates a new migration file
func (m *Migrator) Create(name string) error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := filepath.Join("migrations", filename)

	// Ensure migrations directory exists
	if err := os.MkdirAll("migrations", 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create migration file template
	template := fmt.Sprintf(`-- +migrate Up
-- SQL for migration %s goes here

-- +migrate Down
-- SQL for rollback goes here
`, name)

	if err := ioutil.WriteFile(filepath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	m.logger.Info("Created migration file", zap.String("file", filepath))
	return nil
}

// Status returns the current migration status
func (m *Migrator) Status() (string, error) {
	// Ensure migrations table exists
	if err := m.createMigrationsTable(); err != nil {
		return "", fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return "", fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get migration files: %w", err)
	}

	var status strings.Builder
	status.WriteString("Migration Status:\n")
	status.WriteString("=================\n")

	for _, file := range migrationFiles {
		name := strings.TrimSuffix(file, ".sql")
		if _, exists := applied[name]; exists {
			status.WriteString(fmt.Sprintf("✓ %s (applied)\n", file))
		} else {
			status.WriteString(fmt.Sprintf("✗ %s (pending)\n", file))
		}
	}

	return status.String(), nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// getAppliedMigrations returns a map of applied migrations
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	query := "SELECT migration FROM schema_migrations"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		applied[migration] = true
	}

	return applied, rows.Err()
}

// getMigrationFiles returns sorted migration files
func (m *Migrator) getMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(filename string) error {
	// Read migration file
	content, err := ioutil.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Extract up migration (everything between +migrate Up and +migrate Down)
	lines := strings.Split(string(content), "\n")
	var upSQL strings.Builder
	inUpSection := false

	for _, line := range lines {
		if strings.Contains(line, "+migrate Up") {
			inUpSection = true
			continue
		}
		if strings.Contains(line, "+migrate Down") {
			break
		}
		if inUpSection && !strings.HasPrefix(strings.TrimSpace(line), "--") {
			upSQL.WriteString(line + "\n")
		}
	}

	// Execute migration in a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(upSQL.String()); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	migrationName := strings.TrimSuffix(filename, ".sql")
	if _, err := tx.Exec("INSERT INTO schema_migrations (migration) VALUES ($1)", migrationName); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

// getLastMigration returns the last applied migration
func (m *Migrator) getLastMigration() (string, error) {
	query := "SELECT migration FROM schema_migrations ORDER BY applied_at DESC LIMIT 1"
	var migration string
	err := m.db.QueryRow(query).Scan(&migration)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return migration, err
}

// rollbackMigration rolls back a migration
func (m *Migrator) rollbackMigration(migrationName string) error {
	filename := migrationName + ".sql"
	
	// Read migration file
	content, err := ioutil.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Extract down migration
	lines := strings.Split(string(content), "\n")
	var downSQL strings.Builder
	inDownSection := false

	for _, line := range lines {
		if strings.Contains(line, "+migrate Down") {
			inDownSection = true
			continue
		}
		if inDownSection && !strings.HasPrefix(strings.TrimSpace(line), "--") {
			downSQL.WriteString(line + "\n")
		}
	}

	// Execute rollback in a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute rollback SQL
	if _, err := tx.Exec(downSQL.String()); err != nil {
		return fmt.Errorf("failed to execute rollback SQL: %w", err)
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE migration = $1", migrationName); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return tx.Commit()
}