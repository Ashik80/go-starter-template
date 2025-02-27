package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoDetailsInfoOob struct {
	tmpl *template.Template
}

func NewTodoDetailsInfoOob() TodoDetailsInfoOob {
	return TodoDetailsInfoOob{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (t *TodoDetailsInfoOob) Execute(w http.ResponseWriter, data *entity.Todo) error {
	return t.tmpl.ExecuteTemplate(w, "todo-details-info-oob", data)
}
