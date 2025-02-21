package auth_helpers

import (
	"net/http"
	"unicode"

	"go-starter-template/pkg/entity"
)

func IsStrongPassword(password string) []string {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	errors := []string{}
	if len(password) > 7 {
		hasMinLen = true
	}
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			hasSpecial = true
		}
	}
	if !hasMinLen {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	if !hasUpper {
		errors = append(errors, "Password must have at least 1 uppercase character")
	}
	if !hasLower {
		errors = append(errors, "Password must have at least 1 lowercase character")
	}
	if !hasNumber {
		errors = append(errors, "Password must have at least 1 digit")
	}
	if !hasSpecial {
		errors = append(errors, "Password must have at least 1 special character or symbol")
	}
	return errors
}

func SetSessionCookie(w http.ResponseWriter, session *entity.Session, env string) {
	secure := true
	if env == "development" {
		secure = false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
