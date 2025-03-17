package middlewares

import (
	"net/http"
	"strings"
)

func EnableCors(allowedOrigins string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			origins := strings.Split(allowedOrigins, ",")
			origin := r.Header.Get("Origin")
			for _, o := range origins {
				if origin == o {
					w.Header().Add("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
