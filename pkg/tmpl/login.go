package tmpl

import (
	"html/template"
	"net/http"

	"go-starter-template/pkg/form"
	"go-starter-template/pkg/infrastructure/renderer"
)

type LoginPage struct {
	tmpl *template.Template
}

type LoginPageData struct {
	PageData
	Form *form.LoginForm
}

func NewLoginPage() *LoginPage {
	return &LoginPage{
		tmpl: renderer.GetPageTemplate("login"),
	}
}

func (l *LoginPage) Execute(w http.ResponseWriter, data *form.LoginForm) error {
	pageData := LoginPageData{
		PageData: PageData{
			Title: "Login",
			Path:  "/login",
		},
		Form: data,
	}
	return l.tmpl.ExecuteTemplate(w, "base", pageData)
}
