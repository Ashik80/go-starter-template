package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-starter-template/internal/domain/entities"
	"go-starter-template/internal/domain/repositories"
	"go-starter-template/internal/domain/valueobject"
)

type TodoDTO struct {
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t TodoDTO) toTodo() *entities.Todo {
	return &entities.Todo{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		CreatedAt:   valueobject.NewTime(t.CreatedAt),
		UpdatedAt:   valueobject.NewTime(t.UpdatedAt),
	}
}

type PQTodoRepository struct {
	db *sql.DB
}

func NewPQTodoRepository(db *sql.DB) repositories.ITodoRepository {
	return &PQTodoRepository{db}
}

func (t *PQTodoRepository) List(ctx context.Context) ([]*entities.Todo, error) {
	rows, err := t.db.QueryContext(ctx, "SELECT * FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	todos := make([]*entities.Todo, 0)
	for rows.Next() {
		var todo TodoDTO
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo.toTodo())
	}
	return todos, nil
}

func (t *PQTodoRepository) Get(ctx context.Context, id int) (*entities.Todo, error) {
	row := t.db.QueryRowContext(ctx, "SELECT * FROM todos WHERE id = $1", id)

	var todo TodoDTO
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
	return todo.toTodo(), nil
}

func (t *PQTodoRepository) Create(ctx context.Context, todo *entities.Todo) (*entities.Todo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var createdTodo TodoDTO
	err = tx.QueryRowContext(ctx,
		"INSERT INTO todos (title, description) VALUES ($1, $2) RETURNING *",
		todo.Title,
		todo.Description,
	).Scan(
		&createdTodo.ID,
		&createdTodo.Title,
		&createdTodo.Description,
		&createdTodo.CreatedAt,
		&createdTodo.UpdatedAt,
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
	return createdTodo.toTodo(), nil
}

func (t *PQTodoRepository) Update(ctx context.Context, todo *entities.Todo) (*entities.Todo, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var updatedTodo TodoDTO
	err = tx.QueryRowContext(ctx,
		"UPDATE todos SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING *",
		todo.Title,
		todo.Description,
		todo.ID,
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
	return updatedTodo.toTodo(), nil
}

func (t *PQTodoRepository) Delete(ctx context.Context, todo *entities.Todo) error {
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
