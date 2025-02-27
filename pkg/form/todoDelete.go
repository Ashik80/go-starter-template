package form

import "net/http"

type TodoDeleteForm struct {
	Form
	ID int
}

func NewTodoDeleteForm(r *http.Request, id int) *TodoDeleteForm {
	return &TodoDeleteForm{
		Form: NewForm(r),
		ID:   id,
	}
}
