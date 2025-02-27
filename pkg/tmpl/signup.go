package tmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type SignupPage struct {
	tmpl *template.Template
}

type SignupPageData struct {
	PageData
	Form *form.SignupForm
}

func NewSignupPage() *SignupPage {
	return &SignupPage{
		tmpl: renderer.GetPageTemplate("signup"),
	}
}

func (s *SignupPage) Execute(w http.ResponseWriter, data *form.SignupForm) error {
	pageData := SignupPageData{
		PageData: PageData{
			Title: "Signup",
			Path:  "/signup",
		},
		Form: data,
	}
	return s.tmpl.ExecuteTemplate(w, "base", pageData)
}
