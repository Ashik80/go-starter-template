package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type LoginForm struct {
	tmpl *template.Template
}

func NewLoginForm() *LoginForm {
	return &LoginForm{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (p *LoginForm) Execute(w http.ResponseWriter, data *form.LoginForm) error {
	return p.tmpl.ExecuteTemplate(w, "login-form", data)
}
