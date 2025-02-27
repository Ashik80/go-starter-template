package form

import "net/http"

type SignupForm struct {
	Form
	Email    string
	Password string
	Error    struct {
		Email    string
		Password []string
	}
	FormError string
}

func NewSignupForm(r *http.Request) *SignupForm {
	return &SignupForm{
		Form: NewForm(r),
	}
}
