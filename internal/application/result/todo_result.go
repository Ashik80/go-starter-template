package result

import "go-starter-template/internal/domain/entities"

type TodoResult struct {
	ID          int
	Title       string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

func NewTodoResult(todo *entities.Todo) *TodoResult {
	return &TodoResult{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt.ToString(),
		UpdatedAt:   todo.UpdatedAt.ToString(),
	}
}
