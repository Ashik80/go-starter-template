package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoEditForm struct {
	tmpl *template.Template
}

func NewTodoEditForm() TodoEditForm {
	return TodoEditForm{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (t *TodoEditForm) Execute(w http.ResponseWriter, data *form.TodoEditForm) error {
	return t.tmpl.ExecuteTemplate(w, "todo-edit-form", data)
}
