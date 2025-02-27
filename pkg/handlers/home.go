package handlers

import (
	"net/http"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/tmpl"
)

type HomeHandler struct {
	infrastructure.Router
}

func init() {
	Register(new(HomeHandler))
}

func (h *HomeHandler) Init(a *app.App) error {
	h.Router = a.Router
	return nil
}

func (h *HomeHandler) Routes() {
	h.Get("/", h.Index)
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	p := tmpl.NewHomePage()
	w.WriteHeader(200)
	p.Execute(w, nil)
}
