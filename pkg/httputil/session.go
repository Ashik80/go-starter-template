package httputil

import (
	"net/http"
	"time"

	"go-starter-template/pkg/application/result"
)

const (
	sessionCookieName = "session_id"
)

func SetSessionCookie(w http.ResponseWriter, session *result.SessionResult, env string) {
	secure := env != "development"
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.ID.String(),
		Path:     "/",
		Expires:  session.ExpiresAt,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func RemoveSessionCookie(w http.ResponseWriter, env string) {
	secure := env != "development"
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: secure,
		SameSite: http.SameSiteLaxMode,
	})
}
