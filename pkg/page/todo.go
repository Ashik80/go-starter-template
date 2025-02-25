package page

import (
	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/form"
	"net/http"
)

type TodoPageData struct {
	Form       *TodoEditForm
	Todo       *entity.Todo
	DeleteForm *TodoDeleteForm
}

type TodoEditForm struct {
	form.Form
	ID          int
	Title       string
	Description string
	Error       string
}

type TodoDeleteForm struct {
	form.Form
	ID int
}

func NewTodoPage(r *http.Request) *Page {
	return &Page{
		Title:  "Todo Details",
		Name:   "todo-details",
		Path:   "/todo/" + r.URL.Path,
		Layout: "main",
	}
}

func NewTodoEditForm(r *http.Request, todo *entity.Todo) *TodoEditForm {
	form := &TodoEditForm{
		Form:  form.NewForm(r),
		Error: "",
	}
	
	if todo != nil {
		form.ID = todo.ID
		form.Title = todo.Title
		form.Description = todo.Description
	}
	
	return form
}

func NewTodoPageData(r *http.Request, todo *entity.Todo) *TodoPageData {
	return &TodoPageData{
		Todo:       todo,
		Form:       NewTodoEditForm(r, todo),
		DeleteForm: NewDeleteForm(r, todo),
	}
}

func NewDeleteForm(r *http.Request, todo *entity.Todo) *TodoDeleteForm {
	return &TodoDeleteForm{
		Form: form.NewForm(r),
		ID:   todo.ID,
	}
}
