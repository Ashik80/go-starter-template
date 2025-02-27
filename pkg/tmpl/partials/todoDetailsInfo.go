package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/infrastructure/renderer"
)

type TodoDetailsInfo struct {
	tmpl *template.Template
}

type TodoDetailsInfoData struct {
	*entity.Todo
}

func NewTodoDetailsInfo() TodoDetailsInfo {
	return TodoDetailsInfo{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (t *TodoDetailsInfo) Execute(w http.ResponseWriter, data TodoDetailsInfoData) error {
	return t.tmpl.ExecuteTemplate(w, "todoDetailsInfo", data)
}
