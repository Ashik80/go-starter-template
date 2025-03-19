package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-starter-template/pkg/infrastructure/logger"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	if crw.statusCode == 0 {
		crw.statusCode = 200
	}
	n, err := crw.ResponseWriter.Write(b)
	crw.bytesWritten += int64(n)
	return n, err
}

func Logger(next http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/web/") {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		crw := &CustomResponseWriter{ResponseWriter: w}
		next.ServeHTTP(crw, r)
		duration := time.Since(start)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
		if r.URL.RawQuery != "" {
			url = fmt.Sprintf("%s?%s", url, r.URL.RawQuery)
		}

		logMessage := fmt.Sprintf(
			"\"%s %s %s\" from %s - %d %dB completed in %s",
			r.Method,
			url,
			r.Proto,
			r.RemoteAddr,
			crw.statusCode,
			crw.bytesWritten,
			duration,
		)

		switch {
		case crw.statusCode >= 400:
			log.Error(logMessage)
		case crw.statusCode >= 300:
			log.Info(logMessage)
		default:
			log.Success(logMessage)
		}
	})
}
