package database

import (
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/sirupsen/logrus"
)

type DB struct {
	*sql.DB
	logger *logrus.Logger
}

func New(db *sql.DB) *DB {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	return &DB{
		DB:     db,
		logger: logger,
	}
}

func (db *DB) Logger() *logrus.Logger {
	return db.logger
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				db.logger.Errorf("transaction rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				db.logger.Errorf("transaction commit failed: %v", err)
			}
		}
	}()

	err = fn(tx)
	return err
}

// ExecStatement executes a Jet statement and returns the result
func (db *DB) ExecStatement(stmt postgres.Statement) (sql.Result, error) {
	query, args := stmt.Sql()
	db.logger.Debugf("Executing query: %s with args: %v", query, args)
	return db.Exec(query, args...)
}

// QueryStatement executes a Jet query statement
func (db *DB) QueryStatement(stmt postgres.Statement) (*sql.Rows, error) {
	query, args := stmt.Sql()
	db.logger.Debugf("Executing query: %s with args: %v", query, args)
	return db.Query(query, args...)
}

// QueryRowStatement executes a Jet query statement that returns a single row
func (db *DB) QueryRowStatement(stmt postgres.Statement) *sql.Row {
	query, args := stmt.Sql()
	db.logger.Debugf("Executing query: %s with args: %v", query, args)
	return db.QueryRow(query, args...)
}