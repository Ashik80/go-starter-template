package entities

import (
	"errors"

	"go-starter-template/pkg/domain/valueobject"
)

var (
	ErrTodoIsRequired  = errors.New("todo is required")
	ErrTitleIsRequired = errors.New("title is required")
)

type Todo struct {
	ID          int
	Title       string
	Description string
	CreatedAt   valueobject.Time
	UpdatedAt   valueobject.Time
}

func NewTodo(title, description string) (*Todo, error) {
	if title == "" {
		return nil, ErrTitleIsRequired
	}
	currentTime := valueobject.NewCurrentTime()
	return &Todo{
		Title:       title,
		Description: description,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}, nil
}

func NewTodoWithID(id int, title, description string) (*Todo, error) {
	todo, err := NewTodo(title, description)
	if err != nil {
		return nil, err
	}
	todo.ID = id
	return todo, nil
}

func (t *Todo) SetID(id int) error {
	if t == nil {
		return ErrTodoIsRequired
	}
	t.ID = id
	t.UpdatedAt = valueobject.NewCurrentTime()
	return nil
}

func (t *Todo) UpdateTitle(title string) error {
	if t == nil {
		return ErrTodoIsRequired
	}
	t.Title = title
	t.UpdatedAt = valueobject.NewCurrentTime()
	return nil
}
