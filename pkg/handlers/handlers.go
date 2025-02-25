package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-starter-template/pkg/app"
)

type Handler interface {
	Routes()
	Init(app *app.App) error
}

var handlers []Handler

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}

func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func jsonErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func parseJson(r *http.Request, m any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(m); err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}
	return nil
}

func parseToInt(value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}
	return n, nil
}

func RegisterRoutes(a *app.App) error {
	handlers := GetHandlers()
	for _, h := range handlers {
		if err := h.Init(a); err != nil {
			return fmt.Errorf("failed to initialize handler: %w", err)
		}
		h.Routes()
	}

	return nil
}
