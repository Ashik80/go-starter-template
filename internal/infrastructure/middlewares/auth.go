package middlewares

import (
	"context"
	"net/http"

	"go-starter-template/internal/application/services"
	"go-starter-template/internal/httputil"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

func AuthMiddleware(env string, sessionService services.ISessionService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			result, err := sessionService.GetSession(ctx, cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if result.Session.Expired() {
				sessionService.DeleteSession(ctx, cookie.Value)
				httputil.RemoveSessionCookie(w, env)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx = context.WithValue(ctx, userContextKey, result.Session.User)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(ctx context.Context) any {
	return ctx.Value(userContextKey)
}
