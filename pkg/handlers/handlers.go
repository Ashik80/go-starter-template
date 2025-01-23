package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/service"
)

type Handler interface {
	Routes(service.Router)
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
		return fmt.Errorf("failed to parse json: %v", err)
	}
	return nil
}

func RegisterRoutes(a *app.App) error {
	handlers := GetHandlers()
	for _, h := range handlers {
		if err := h.Init(a); err != nil {
			return fmt.Errorf("failed to initialize handler: %v", err)
		}
		h.Routes(a.Router)
	}
	return nil
}
