package controllers

import (
	"net/http"

	"go-starter-template/pkg/infrastructure/router"
	"go-starter-template/pkg/infrastructure/views/pages"
)

type HomeController struct{}

func NewHomeController(r router.Router) {
	controller := &HomeController{}
	r.Get("/", controller.Index)
}

func (hc *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(r.Context(), w)
}
