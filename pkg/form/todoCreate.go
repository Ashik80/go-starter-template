package form

import "net/http"

type TodoCreateForm struct {
	Form
	Title       string
	Description string
	Error       string
}

func NewTodoCreateForm(r *http.Request) *TodoCreateForm {
	return &TodoCreateForm{
		Form: NewForm(r),
	}
}
