package form

import (
	"net/http"

	"go-starter-template/pkg/entity"
)

type TodoEditForm struct {
	Form
	*entity.Todo
	Error string
}

func NewTodoEditForm(r *http.Request, todo *entity.Todo) *TodoEditForm {
	form := &TodoEditForm{
		Form: NewForm(r),
		Todo: todo,
	}

	return form
}
