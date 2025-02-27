package tmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodosPage struct {
	template *template.Template
}

type TodosPageData struct {
	*PageData
	Form  *form.TodoCreateForm
	Todos []*entity.Todo
	Error string
}

func NewTodosPage() *TodosPage {
	return &TodosPage{
		template: renderer.GetPageTemplate("todos"),
	}
}

func (p *TodosPage) Execute(w http.ResponseWriter, data *TodosPageData) error {
	data.PageData = &PageData{
		Title: "Todos",
		Path:  "/todos",
	}
	return p.template.ExecuteTemplate(w, "base", data)
}
