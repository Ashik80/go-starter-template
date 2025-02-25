package middlewares

import (
	"context"
	"net/http"
	"time"

	"go-starter-template/pkg/helpers/auth_helpers"
	"go-starter-template/pkg/service"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

func AuthMiddleware(env string, sessionService service.SessionService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			s, err := sessionService.GetSession(ctx, cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if s.ExpiresAt.Unix() < time.Now().Unix() {
				sessionService.DeleteSession(ctx, cookie.Value)
				auth_helpers.RemoveSessionCookie(w, env)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx = context.WithValue(ctx, userContextKey, s.User)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(ctx context.Context) any {
	return ctx.Value(userContextKey)
}
