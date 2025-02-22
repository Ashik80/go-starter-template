package handlers

import (
	"net/http"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/page"
	"go-starter-template/pkg/service"
)

type HomeHandler struct {
	service.TemplateRenderer
	service.Router
}

func init() {
	Register(new(HomeHandler))
}

func (h *HomeHandler) Init(a *app.App) error {
	h.Router = a.Router
	h.TemplateRenderer = a.TemplateRenderer
	return nil
}

func (h *HomeHandler) Routes() {
	h.Router.HandleFunc("/", h.Index)
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	p := page.NewHomePage()
	w.WriteHeader(200)
	h.Render(w, p)
}
