package partialTmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type SignupForm struct {
	tmpl *template.Template
}

func NewSignupForm() *SignupForm {
	return &SignupForm{
		tmpl: renderer.GetBaseTemplate(),
	}
}

func (p *SignupForm) Execute(w http.ResponseWriter, data *form.SignupForm) error {
	return p.tmpl.ExecuteTemplate(w, "signup-form", data)
}
