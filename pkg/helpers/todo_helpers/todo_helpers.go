package todo_helpers

import (
	"go-starter-template/pkg/page"
	"strings"
)

func ValidateTodoForm(form *page.TodoCreateForm) (bool, []string) {
	if strings.TrimSpace(form.Title) == "" {
		return false, []string{"Title cannot be empty"}
	}
	return true, nil
}
