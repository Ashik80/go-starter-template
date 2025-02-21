package page

import (
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/store"
	"net/http"
)

type TodoListPageData struct {
	Todos []*store.Todo
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

func NewTodoForm(r *http.Request) *TodoCreateForm {
	return &TodoCreateForm{
		Form:        form.NewForm(r),
		Title:       "",
		Description: "",
		Error:       "",
	}
}

func NewTodoListPageData(r *http.Request) *TodoListPageData {
	return &TodoListPageData{
		Todos: []*store.Todo{},
		Form:  NewTodoForm(r),
	}
}
