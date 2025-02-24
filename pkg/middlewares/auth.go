package middlewares

import (
	"context"
	"net/http"
	"time"

	"go-starter-template/pkg/store"
)

func AuthMiddleware(sessionStore store.SessionStore) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			s, err := sessionStore.GetWithUser(ctx, cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if s.ExpiresAt.Unix() < time.Now().Unix() {
				sessionStore.Delete(ctx, s.ID.String())
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx = context.WithValue(ctx, "user", s.User)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
