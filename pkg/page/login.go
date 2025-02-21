package page

import (
	"go-starter-template/pkg/form"
	"net/http"
)

type LoginPageData struct {
	Form *LoginForm
}

type LoginForm struct {
	form.Form
	Email    string
	Password string
	Remember string
	Error    string
}

func NewLoginPage() *Page {
	p := New()
	p.Title = "Login"
	p.Path = "/login"
	p.Layout = "auth"
	p.Name = "login"
	return p
}

func NewLoginForm(r *http.Request) *LoginForm {
	return &LoginForm{
		Form:     form.NewForm(r),
		Email:    "",
		Password: "",
		Remember: "",
	}
}

func NewLoginPageData(r *http.Request) *LoginPageData {
	return &LoginPageData{
		Form: NewLoginForm(r),
	}
}
