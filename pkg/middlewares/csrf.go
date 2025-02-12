package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	c := csrf.Protect([]byte("somehting-in-the-way"))
	return c(next)
}
