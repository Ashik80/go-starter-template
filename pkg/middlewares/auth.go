package middlewares

import (
	"context"
	"net/http"

	"go-starter-template/pkg/store"
)

type MiddlewareHandler func(next http.Handler) http.Handler

type Middleware struct {
	sessionStore store.SessionStore
}

func NewMiddleware(sessionStore store.SessionStore) *Middleware {
	return &Middleware{sessionStore}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		s, err := m.sessionStore.GetWithUser(ctx, cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx = context.WithValue(ctx, "user", s.User)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
