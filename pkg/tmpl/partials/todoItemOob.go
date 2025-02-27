package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoItemOob struct {
	tmpl *template.Template
}

func NewTodoItemOob() *TodoItemOob {
	return &TodoItemOob{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (p *TodoItemOob) Execute(w http.ResponseWriter, data *entity.Todo) error {
	return p.tmpl.ExecuteTemplate(w, "todo-item-oob", data)
}
