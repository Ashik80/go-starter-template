package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/form"
	"go-starter-template/pkg/helpers/auth_helpers"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"

	"golang.org/x/crypto/bcrypt"
)

type (
	AuthHandlers struct {
		*service.TemplateRenderer
		service.Router
		env          string
		userStore    store.UserStore
		sessionStore store.SessionStore
	}

	LoginPageData struct {
		Form LoginForm
	}

	SignupPageData struct {
		Form SignupForm
	}

	LoginForm struct {
		form.Form
		Email    string
		Password string
		Remember string
		Error    string
	}

	SignupForm struct {
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
)

func newLoginPage() *page.Page {
	p := page.New()
	p.Layout = "auth"
	p.Name = "login"
	return p
}

func newSignupPage() *page.Page {
	p := page.New()
	p.Layout = "auth"
	p.Name = "signup"
	return p
}

func newLoginForm(r *http.Request) LoginForm {
	return LoginForm{
		Email:    "",
		Password: "",
		Remember: "",
		Form:     form.NewForm(r),
	}
}

func newSignupForm(r *http.Request) SignupForm {
	return SignupForm{
		Email:    "",
		Password: "",
		Form:     form.NewForm(r),
	}
}

func init() {
	Register(new(AuthHandlers))
}

func (h *AuthHandlers) Init(a *app.App) error {
	h.Router = a.Router
	h.userStore = a.Store.UserStore
	h.env = a.Config.Env
	h.sessionStore = a.Store.SessionStore
	h.TemplateRenderer = a.TemplateRenderer
	return nil
}

func (h *AuthHandlers) Routes() {
	h.HandleFunc("/login", h.LoginView)
	h.HandleFunc("POST /login", h.Login)
	h.HandleFunc("/signup", h.SignupView)
	h.HandleFunc("POST /signup", h.Signup)
}

func (h *AuthHandlers) LoginView(w http.ResponseWriter, r *http.Request) {
	p := newLoginPage()
	loginPageData := &LoginPageData{
		Form: newLoginForm(r),
	}
	p.Data = loginPageData
	h.Render(w, p)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := newLoginForm(r)
	loginForm.Email = r.FormValue("email")
	loginForm.Password = r.FormValue("password")
	loginForm.Remember = r.FormValue("remember")

	ctx := r.Context()

	user, err := h.userStore.GetByEmail(ctx, loginForm.Email)
	if err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		h.RenderPartial(w, "login-form", loginForm)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		loginForm.Error = "Invalid credentials"
		w.WriteHeader(401)
		h.RenderPartial(w, "login-form", loginForm)
		return
	}

	var sessionExpiry time.Time
	if loginForm.Remember == "on" {
		sessionExpiry = time.Now().Add(30 * 24 * time.Hour)
	} else {
		sessionExpiry = time.Now().Add(1 * time.Hour)
	}

	sess, err := h.sessionStore.Create(ctx, user, sessionExpiry)
	if err != nil {
		log.Printf("ERROR: %v", err)
		loginForm.Error = "Something went wrong"
		w.WriteHeader(500)
		h.RenderPartial(w, "login-form", loginForm)
	}

	auth_helpers.SetSessionCookie(w, sess, h.env)

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}

func (h *AuthHandlers) SignupView(w http.ResponseWriter, r *http.Request) {
	p := newSignupPage()
	signupPageData := &SignupPageData{
		Form: newSignupForm(r),
	}
	p.Data = signupPageData
	h.Render(w, p)
}

func (h *AuthHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signupForm := newSignupForm(r)
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

	user, err := h.userStore.GetByEmail(ctx, signupForm.Email)
	if err != nil {
		var notFoundError *store.NotFoundError
		if !errors.As(err, &notFoundError) {
			log.Printf("ERROR: %v", err)
			signupForm.Error.ErrorMessage = err.Error()
			w.WriteHeader(500)
			h.RenderPartial(w, "signup-form", signupForm)
			return
		}
	}

	if user != nil {
		errorMsg := fmt.Sprintf("user with email %s already exists", signupForm.Email)
		log.Printf("ERROR: %s", errorMsg)
		signupForm.Error.ErrorMessage = errorMsg
		w.WriteHeader(http.StatusConflict)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), 10)
	if err != nil {
		log.Printf("ERROR: failed to generate password hash: %v\n", err)
		signupForm.Error.ErrorMessage = "failed to generate password hash"
		w.WriteHeader(500)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	_, err = h.userStore.Create(ctx, signupForm.Email, string(passwordHash))
	if err != nil {
		signupForm.Error.ErrorMessage = err.Error()
		log.Printf("ERROR: failed to create user: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.RenderPartial(w, "signup-form", signupForm)
		return
	}

	w.WriteHeader(200)
	h.RenderPartial(w, "signup-form", newSignupForm(r))
	h.RenderPartial(w, "signup-success-message", nil)
}
