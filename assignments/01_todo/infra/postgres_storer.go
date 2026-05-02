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

func (ps *PostgresStore) Create(task todo.Task) (todo.Task, error) {
	err := ps.db.QueryRow(
		`INSERT INTO tasks (title, description, priority, done, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		task.Title,
		task.Description,
		task.Priority,
		task.Done,
		task.CreatedAt,
	).Scan(&task.Id, &task.CreatedAt)

	if err != nil {
		return todo.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (ps *PostgresStore) List(filter string) ([]todo.Task, error) {
	query := `SELECT id, title, description, priority, done, created_at FROM tasks`
	switch filter {
	case "done":
		query += ` WHERE done = true`
	case "pending":
		query += ` WHERE done = false`
	}
	query += ` ORDER BY id`

	rows, err := ps.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list task: %w", err)
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

func (ps *PostgresStore) GetByID(id int) (todo.Task, error) {
	var task todo.Task
	err := ps.db.QueryRow(
		`SELECT id, title, description, priority, done, created_at
		FROM tasks
		WHERE id = $1`,
		id,
	).Scan(&task.Id, &task.Title, &task.Description, &task.Priority, &task.Done, &task.CreatedAt)

	if err == sql.ErrNoRows {
		return todo.Task{}, fmt.Errorf("id not found")
	} else if err != nil {
		return todo.Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}

// Update modifies an existing task by ID and returns the updated task.
func (ps *PostgresStore) Update(task todo.Task) (todo.Task, error) {
	err := ps.db.QueryRow(
		`UPDATE tasks
		 SET title = $1, description = $2, priority = $3
		 WHERE id = $4
		 RETURNING id, title, description, priority, done, created_at`,
		task.Title, task.Description, task.Priority, task.Id,
	).Scan(&task.Id, &task.Title, &task.Description, &task.Priority, &task.Done, &task.CreatedAt)

	if err == sql.ErrNoRows {
		return todo.Task{}, fmt.Errorf("task #%d not found", task.Id)
	} else if err != nil {
		return todo.Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

// Delete removes a task by ID.
func (ps *PostgresStore) Delete(id int) error {
	result, err := ps.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task #%d not found", id)
	}

	return nil
}

// ToggleDone flips the done status of a task by ID and returns the updated task.
func (ps *PostgresStore) ToggleDone(id int) (todo.Task, error) {
	var task todo.Task
	err := ps.db.QueryRow(
		`UPDATE tasks
		 SET done = NOT done
		 WHERE id = $1
		 RETURNING id, title, description, priority, done, created_at`,
		id,
	).Scan(&task.Id, &task.Title, &task.Description, &task.Priority, &task.Done, &task.CreatedAt)

	if err == sql.ErrNoRows {
		return todo.Task{}, fmt.Errorf("task #%d not found", id)
	} else if err != nil {
		return todo.Task{}, fmt.Errorf("toggle done: %w", err)
	}

	return task, nil
}
