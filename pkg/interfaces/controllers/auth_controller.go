package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/domain/valueobject"
	"go-starter-template/pkg/httputil"
	"go-starter-template/pkg/infrastructure/config"
	"go-starter-template/pkg/infrastructure/csrf"
	"go-starter-template/pkg/infrastructure/router"
	"go-starter-template/pkg/infrastructure/views/components"
	"go-starter-template/pkg/infrastructure/views/pages"
)

type AuthController struct {
	userService interfaces.UserService
	config      *config.Config
}

type LoginForm struct {
	CSRF     template.HTML
	Email    string
	Password string
	Remember string
	Error    string
}

type SignupForm struct {
	CSRF     template.HTML
	Email    string
	Password string
	Error    struct {
		Email    string
		Password []string
	}
	FormError string
}

func NewAuthController(r router.Router, userService interfaces.UserService, config *config.Config) {
	controller := &AuthController{
		userService: userService,
		config:      config,
	}

	r.Get("/login", controller.LoginView)
	r.Post("/login", controller.Login)
	r.Get("/signup", controller.SignupView)
	r.Post("/signup", controller.Signup)
}

func (ac *AuthController) LoginView(w http.ResponseWriter, r *http.Request) {
	form := components.NewLoginFormData(r)
	w.WriteHeader(200)
	pages.Login(form).Render(r.Context(), w)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	form := &components.LoginFormData{
		CSRF:     csrf.GetCSRFField(r),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
		Remember: r.FormValue("remember"),
	}

	result, err := ac.userService.Login(r.Context(), &command.CreateLoginCommand{
		Email:    form.Email,
		Password: form.Password,
		Remember: form.Remember == "on",
	})

	if err != nil {
		form.Error = err.Error()
		w.WriteHeader(401)
		components.LoginForm(form).Render(r.Context(), w)
		return
	}

	httputil.SetSessionCookie(w, result.Session, ac.config.Env)

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}

func (ac *AuthController) SignupView(w http.ResponseWriter, r *http.Request) {
	form := components.NewSignupFormData(r)
	w.WriteHeader(200)
	pages.Signup(form).Render(r.Context(), w)
}

func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	form := &components.SignupFormData{
		CSRF:     csrf.GetCSRFField(r),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	_, err := ac.userService.Signup(r.Context(), &command.CreateSignupCommand{
		Email:    form.Email,
		Password: form.Password,
	})

	if err != nil {
		if err == valueobject.ErrEmailIsInvalid {
			form.Error.Email = err.Error()
		} else {
			form.FormError = err.Error()
		}
		var pve *valueobject.PasswordValidationError
		if errors.As(err, &pve) {
			form.Error.Password = pve.Errors
		}
		w.WriteHeader(400)
		components.SignupForm(form).Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Location", "/login")
	w.WriteHeader(200)
}
