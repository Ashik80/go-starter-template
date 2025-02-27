package tmpl

import (
	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
	"html/template"
	"net/http"
	"strconv"
)

type TodoDetailsPage struct {
	template *template.Template
}

type TodoDetailsPageData struct {
	*PageData
	Todo       *entity.Todo
	EditForm   *form.TodoEditForm
	DeleteForm *form.TodoDeleteForm
	Error      string
}

func NewTodoDetailsPage() TodoDetailsPage {
	return TodoDetailsPage{
		template: renderer.GetPageTemplate("todo-details"),
	}
}

func (t *TodoDetailsPage) GetPageTemplate() *template.Template {
	return t.template
}

func (t *TodoDetailsPage) Execute(w http.ResponseWriter, data *TodoDetailsPageData) error {
	id := strconv.Itoa(data.Todo.ID)
	data.PageData = &PageData{
		Title: data.Todo.Title,
		Path:  "/todos/" + id,
	}
	return t.template.ExecuteTemplate(w, "base", data)
}
