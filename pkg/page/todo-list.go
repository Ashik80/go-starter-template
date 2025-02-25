package page

import (
	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/form"
	"net/http"
)

type TodoListPageData struct {
	Todos []*entity.Todo
	Form  *TodoCreateForm
}

type TodoCreateForm struct {
	form.Form
	Title       string
	Description string
	Error       string
}

func NewTodoListPage() *Page {
	return &Page{
		Name:   "todo",
		Title:  "Todos",
		Layout: "main",
		Path:   "/todos",
	}
}

func NewTodoCreateForm(r *http.Request) *TodoCreateForm {
	return &TodoCreateForm{
		Form:  form.NewForm(r),
		Error: "",
	}
}

func NewTodoListPageData(r *http.Request) *TodoListPageData {
	return &TodoListPageData{
		Todos: []*entity.Todo{},
		Form:  NewTodoCreateForm(r),
	}
}
