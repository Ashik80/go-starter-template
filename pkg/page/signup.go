package page

import (
	"go-starter-template/pkg/form"
	"net/http"
)

type SignupPageData struct {
	Form *SignupForm
}

type SignupForm struct {
	form.Form
	Email    string
	Password string
	Error    struct {
		Email    string
		Password struct {
			Validations []string
		}
		ErrorMessage string
	}
}

func NewSignupPage() *Page {
	p := New()
	p.Title = "Signup"
	p.Name = "signup"
	p.Layout = "auth"
	p.Path = "/signup"
	return p
}

func NewSignupPageData(r *http.Request) *SignupPageData {
	return &SignupPageData{
		Form: NewSignupForm(r),
	}
}

func NewSignupForm(r *http.Request) *SignupForm {
	return &SignupForm{
		Form:     form.NewForm(r),
		Email:    "",
		Password: "",
	}
}
