package form

import "net/http"

type LoginForm struct {
	Form
	Email    string
	Password string
	Remember string
	Error    string
}

func NewLoginForm(r *http.Request) *LoginForm {
	return &LoginForm{
		Form: NewForm(r),
	}
}
