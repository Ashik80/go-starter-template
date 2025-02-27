package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-starter-template/pkg/entity"
)

type TodoCreateDto struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type PQTodoStore struct {
	db *sql.DB
}

func NewPQTodoStore(db *sql.DB) TodoRepository {
	return &PQTodoStore{db}
}

func (t *PQTodoStore) List(ctx context.Context) ([]*entity.Todo, error) {
	rows, err := t.db.QueryContext(ctx, "SELECT * FROM todos")
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	todos := make([]*entity.Todo, 0)
	for rows.Next() {
		var todo entity.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (t *PQTodoStore) Get(ctx context.Context, id int) (*entity.Todo, error) {
	row := t.db.QueryRowContext(ctx, "SELECT * FROM todos WHERE id = $1", id)

	var todo entity.Todo
	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan todo: %w", err)
	}
	return &todo, nil
}

func (t *PQTodoStore) Create(ctx context.Context, todoDto *TodoCreateDto) (*entity.Todo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var todo entity.Todo
	err = tx.QueryRowContext(ctx,
		"INSERT INTO todos (title, description) VALUES ($1, $2) RETURNING *",
		todoDto.Title,
		todoDto.Description,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &todo, nil
}

func (t *PQTodoStore) Update(ctx context.Context, id int, todoDto *TodoCreateDto) (*entity.Todo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var updatedTodo entity.Todo
	err = tx.QueryRowContext(ctx,
		"UPDATE todos SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING *",
		todoDto.Title,
		todoDto.Description,
		id,
	).Scan(
		&updatedTodo.ID,
		&updatedTodo.Title,
		&updatedTodo.Description,
		&updatedTodo.CreatedAt,
		&updatedTodo.UpdatedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &updatedTodo, nil
}

func (t *PQTodoStore) Delete(ctx context.Context, todo *entity.Todo) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM todos WHERE id = $1", todo.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
