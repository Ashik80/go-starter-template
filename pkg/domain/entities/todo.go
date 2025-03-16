package entities

import (
	"errors"
	"go-starter-template/pkg/domain/valueobject"
	"time"
)

var (
	ErrTitleIsRequired = errors.New("title is required")
)

type Todo struct {
	ID          int
	Title       string
	Description string
	CreatedAt   valueobject.Time
	UpdatedAt   valueobject.Time
}

func NewTodo(title, description string) *Todo {
	return &Todo{
		Title:       title,
		Description: description,
		CreatedAt:   valueobject.Time(time.Now()),
		UpdatedAt:   valueobject.Time(time.Now()),
	}
}

func (t *Todo) Validate() error {
	if t.Title == "" {
		return ErrTitleIsRequired
	}
	return nil
}

func (t *Todo) UpdateTitle(title string) error {
	t.Title = title
	t.UpdatedAt = valueobject.Time(time.Now())

	return t.Validate()
}
