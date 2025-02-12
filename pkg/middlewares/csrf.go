package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRFMiddleware(csrfAuthKey string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		c := csrf.Protect([]byte(csrfAuthKey))
		return c(next)
	}
}
