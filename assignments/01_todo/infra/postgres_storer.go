package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
)

// PostgresStore implements Storer using a PostgreSQL database.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgresStore and verifies the connection.
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// Close closes the underlying database connection.
func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}

// Load retrieves all tasks from the database.
func (ps *PostgresStore) Load() ([]todo.Task, error) {
	rows, err := ps.db.Query(
		`SELECT id, title, description, priority, done, created_at
		 FROM tasks
		 ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []todo.Task
	for rows.Next() {
		var t todo.Task
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Priority, &t.Done, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	if tasks == nil {
		tasks = []todo.Task{}
	}

	return tasks, nil
}

// Save replaces all tasks in the database with the provided slice.
// It uses a transaction to delete all existing rows and insert the new ones.
func (ps *PostgresStore) Save(tasks []todo.Task) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM tasks`); err != nil {
		return fmt.Errorf("delete tasks: %w", err)
	}

	for _, t := range tasks {
		_, err := tx.Exec(
			`INSERT INTO tasks (id, title, description, priority, done, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			t.Id, t.Title, t.Description, t.Priority, t.Done, t.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert task %d: %w", t.Id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
