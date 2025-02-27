package tmpl

import (
	"go-starter-template/pkg/infrastructure/renderer"
	"html/template"
	"net/http"
)

type HomePage struct {
	template *template.Template
}

func NewHomePage() *HomePage {
	return &HomePage{
		template: renderer.GetPageTemplate("home"),
	}
}

func (p *HomePage) Execute(w http.ResponseWriter, data any) error {
	data = &PageData{
		Title: "Home",
		Path:  "/",
	}
	return p.template.ExecuteTemplate(w, "base", data)
}
