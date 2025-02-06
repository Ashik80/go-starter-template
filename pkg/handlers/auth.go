package handlers

import (
	"log"
	"net/http"
	"net/mail"
	"time"

	"go-starter-template/pkg/app"
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
		Email    string
		Password string
		Remember string
		Error    string
	}

	SignupForm struct {
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

func newLoginForm() LoginForm {
	return LoginForm{
		Email:    "",
		Password: "",
		Remember: "false",
	}
}

func newSignupForm() SignupForm {
	return SignupForm{
		Email:    "",
		Password: "",
	}
}

func init() {
	Register(new(AuthHandlers))
}

func (h *AuthHandlers) Init(a *app.App) error {
	h.Router = a.Router
	h.userStore = a.Store.UserStore
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
		Form: newLoginForm(),
	}
	p.Data = loginPageData
	h.Render(w, p)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	loginForm := LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Remember: r.FormValue("remember"),
	}

	ctx := r.Context()

	user, err := h.userStore.GetByEmail(ctx, loginForm.Email)
	if err != nil {
		loginForm.Error = "Invalid credentials"
		h.RenderPartial(w, 401, "login-form", loginForm)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		loginForm.Error = "Invalid credentials"
		h.RenderPartial(w, 401, "login-form", loginForm)
		return
	}

	sessionExpiry := time.Now().Add(24 * time.Hour)
	sess, err := h.sessionStore.Create(ctx, user, sessionExpiry)
	if err != nil {
		log.Printf("ERROR: %v", err)
		loginForm.Error = "Something went wrong"
		h.RenderPartial(w, 500, "login-form", loginForm)
	}

	auth_helpers.SetSessionCookie(w, sess)

	w.Header().Add("Hx-Location", "/todos")
	w.WriteHeader(200)
}

func (h *AuthHandlers) SignupView(w http.ResponseWriter, r *http.Request) {
	p := newSignupPage()
	signupPageData := &SignupPageData{
		Form: newSignupForm(),
	}
	p.Data = signupPageData
	h.Render(w, p)
}

func (h *AuthHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	signupForm := SignupForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	_, err := mail.ParseAddress(signupForm.Email)
	if err != nil {
		signupForm.Error.Email = "Invalid email"
		h.RenderPartial(w, 400, "signup-form", signupForm)
		return
	}

	validationErrors := auth_helpers.IsStrongPassword(signupForm.Password)
	if len(validationErrors) > 0 {
		signupForm.Error.Password.Validations = validationErrors
		signupForm.Error.ErrorMessage = "Password must have minimum lenght of 8 characters and must contain at least 1 uppercase letter, 1 lowercase letter, 1 digit and 1 symbol"
		h.RenderPartial(w, 400, "signup-form", signupForm)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), 10)
	if err != nil {
		log.Printf("ERROR: failed to generate password hash: %v\n", err)
		signupForm.Error.ErrorMessage = "failed to generate password hash"
		h.RenderPartial(w, 500, "signup-form", signupForm)
		return
	}

	_, err = h.userStore.Create(r.Context(), signupForm.Email, string(passwordHash))
	if err != nil {
		signupForm.Error.ErrorMessage = err.Error()
		log.Printf("failed to create user: %v", err)
		h.RenderPartial(w, 422, "signup-form", signupForm)
		return
	}

	h.RenderPartial(w, 200, "signup-form", newSignupForm())
	h.RenderPartial(w, 200, "signup-success-message", nil)
}
