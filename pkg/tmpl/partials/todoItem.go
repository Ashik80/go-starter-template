package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoItem struct {
	tmpl *template.Template
}

func NewTodoItem() *TodoItem {
	return &TodoItem{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (p *TodoItem) Execute(w http.ResponseWriter, data *entity.Todo) error {
	return p.tmpl.ExecuteTemplate(w, "todo-item", data)
}
