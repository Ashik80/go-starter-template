package controllers

import (
	"net/http"

	"go-starter-template/internal/infrastructure/views/pages"
	"go-starter-template/pkg/router"
)

type HomeController struct{}

func NewHomeController(r router.Router) {
	controller := &HomeController{}
	r.Get("/", controller.Index)
}

func (hc *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(r.Context(), w)
}
