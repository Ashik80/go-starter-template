package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/helpers/auth_helpers"
	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/repository"
	"go-starter-template/pkg/service"
)

type AuthHandlers struct {
	infrastructure.TemplateRenderer
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
	h.TemplateRenderer = a.TemplateRenderer
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
	p := page.NewLoginPage()
	p.Data = page.NewLoginPageData(r)
	w.WriteHeader(200)
	h.Render(w, p)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := page.NewLoginForm(r)
	loginForm.Email = r.FormValue("email")
	loginForm.Password = r.FormValue("password")
	loginForm.Remember = r.FormValue("remember")

	ctx := r.Context()

	user, err := h.userService.GetUserByEmail(ctx, loginForm.Email)
	if err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		h.RenderPartial(w, "login-form", loginForm)
		return
	}

	if err = h.passwordHasher.CompareHashAndPassword(user.Password, loginForm.Password); err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		h.RenderPartial(w, "login-form", loginForm)
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
		h.RenderPartial(w, "login-form", loginForm)
		return
	}

	auth_helpers.SetSessionCookie(w, sess, h.env)

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}

func (h *AuthHandlers) SignupView(w http.ResponseWriter, r *http.Request) {
	p := page.NewSignupPage()
	p.Data = page.NewSignupPageData(r)
	w.WriteHeader(200)
	h.Render(w, p)
}

func (h *AuthHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signupForm := page.NewSignupForm(r)
	signupForm.Email = r.FormValue("email")
	signupForm.Password = r.FormValue("password")

	_, err := mail.ParseAddress(signupForm.Email)
	if err != nil {
		signupForm.Error.Email = "Invalid email"
		w.WriteHeader(400)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	validationErrors := auth_helpers.IsStrongPassword(signupForm.Password)
	if len(validationErrors) > 0 {
		signupForm.Error.Password.Validations = validationErrors
		signupForm.Error.ErrorMessage = "Password must have minimum lenght of 8 characters and must contain at least 1 uppercase letter, 1 lowercase letter, 1 digit and 1 symbol"
		w.WriteHeader(400)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, signupForm.Email)
	if err != nil {
		var notFoundError *repository.NotFoundError
		if !errors.As(err, &notFoundError) {
			signupForm.Error.ErrorMessage = err.Error()
			w.WriteHeader(500)
			h.RenderPartial(w, "signup-form", signupForm)
			return
		}
	}

	if user != nil {
		signupForm.Error.ErrorMessage = fmt.Sprintf("user with email %s already exists", signupForm.Email)
		w.WriteHeader(http.StatusConflict)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	passwordHash, err := h.passwordHasher.GenerateFromPassword(signupForm.Password, 10)
	if err != nil {
		signupForm.Error.ErrorMessage = "failed to generate password hash"
		w.WriteHeader(500)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	input := service.CreateUserInput{
		Email:    signupForm.Email,
		Password: passwordHash,
	}

	_, err = h.userService.CreateUser(ctx, input)
	if err != nil {
		signupForm.Error.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	w.Header().Add("Hx-Location", "/login")
	w.WriteHeader(200)
}
