package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoCreateForm struct {
	tmpl *template.Template
}

func NewTodoCreateForm() *TodoCreateForm {
	return &TodoCreateForm{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (t *TodoCreateForm) Execute(w http.ResponseWriter, data *form.TodoCreateForm) error {
	return t.tmpl.ExecuteTemplate(w, "todo-create-form", data)
}
