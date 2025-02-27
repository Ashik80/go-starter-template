package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/helpers/auth_helpers"
	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/tmpl"
	partialTmpl "go-starter-template/pkg/tmpl/partials"
)

type AuthHandlers struct {
	infrastructure.Router
	env            string
	userService    service.UserService
	sessionService service.SessionService
	passwordHasher infrastructure.PasswordHasher
}

func init() {
	Register(new(AuthHandlers))
}

func (h *AuthHandlers) Init(a *app.App) error {
	h.Router = a.Router
	h.env = a.Config.Env
	h.userService = a.Services.User
	h.sessionService = a.Services.Session
	h.passwordHasher = a.PasswordHasher
	return nil
}

func (h *AuthHandlers) Routes() {
	h.Get("/login", h.LoginView)
	h.Post("/login", h.Login)
	h.Get("/signup", h.SignupView)
	h.Post("/signup", h.Signup)
}

func (h *AuthHandlers) LoginView(w http.ResponseWriter, r *http.Request) {
	p := tmpl.NewLoginPage()
	loginForm := form.NewLoginForm(r)
	w.WriteHeader(200)
	p.Execute(w, loginForm)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginForm := form.NewLoginForm(r)
	loginForm.Email = strings.TrimSpace(r.FormValue("email"))
	loginForm.Password = strings.TrimSpace(r.FormValue("password"))
	loginForm.Remember = strings.TrimSpace(r.FormValue("remember"))

	loginFormTmpl := partialTmpl.NewLoginForm()

	user, err := h.userService.GetUserByEmail(ctx, loginForm.Email)
	if err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		loginFormTmpl.Execute(w, loginForm)
		return
	}

	if err = h.passwordHasher.CompareHashAndPassword(user.Password, loginForm.Password); err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		loginFormTmpl.Execute(w, loginForm)
		return
	}

	sessionExpiry := time.Now().Add(1 * time.Hour)
	if loginForm.Remember == "on" {
		sessionExpiry = time.Now().Add(30 * 24 * time.Hour)
	}

	sess, err := h.sessionService.CreateSession(ctx, user.ID, sessionExpiry)
	if err != nil {
		loginForm.Error = err.Error()
		w.WriteHeader(500)
		loginFormTmpl.Execute(w, loginForm)
		return
	}

	auth_helpers.SetSessionCookie(w, sess, h.env)

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}

func (h *AuthHandlers) SignupView(w http.ResponseWriter, r *http.Request) {
	p := tmpl.NewSignupPage()
	signupForm := form.NewSignupForm(r)
	w.WriteHeader(200)
	p.Execute(w, signupForm)
}

func (h *AuthHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	signupForm := form.NewSignupForm(r)
	signupForm.Email = strings.TrimSpace(r.FormValue("email"))
	signupForm.Password = strings.TrimSpace(r.FormValue("password"))

	signupFormTmpl := partialTmpl.NewSignupForm()

	_, err := mail.ParseAddress(signupForm.Email)
	if err != nil {
		signupForm.Error.Email = "Invalid email"
		w.WriteHeader(400)
		signupFormTmpl.Execute(w, signupForm)
		return
	}

	validationErrors := auth_helpers.IsStrongPassword(signupForm.Password)
	if len(validationErrors) > 0 {
		signupForm.Error.Password = validationErrors
		signupForm.FormError = "Password must have minimum lenght of 8 characters and must contain at least 1 uppercase letter, 1 lowercase letter, 1 digit and 1 symbol"
		w.WriteHeader(400)
		signupFormTmpl.Execute(w, signupForm)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, signupForm.Email)
	if err != nil {
		err = errors.Unwrap(err)
		if err != sql.ErrNoRows {
			signupForm.FormError = err.Error()
			w.WriteHeader(500)
			signupFormTmpl.Execute(w, signupForm)
			return
		}
	}

	if user != nil {
		signupForm.FormError = fmt.Sprintf("user with email %s already exists", signupForm.Email)
		w.WriteHeader(http.StatusConflict)
		signupFormTmpl.Execute(w, signupForm)
		return
	}

	passwordHash, err := h.passwordHasher.GenerateFromPassword(signupForm.Password, 10)
	if err != nil {
		signupForm.FormError = "failed to generate password hash"
		w.WriteHeader(500)
		signupFormTmpl.Execute(w, signupForm)
		return
	}

	input := service.CreateUserInput{
		Email:    signupForm.Email,
		Password: passwordHash,
	}

	_, err = h.userService.CreateUser(ctx, input)
	if err != nil {
		signupForm.FormError = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		signupFormTmpl.Execute(w, signupForm)
		return
	}

	w.Header().Add("Hx-Location", "/login")
	w.WriteHeader(200)
}
