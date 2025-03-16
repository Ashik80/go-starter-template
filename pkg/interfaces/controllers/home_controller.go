package controllers

import (
	"net/http"

	"go-starter-template/pkg/infrastructure/renderer"
	"go-starter-template/pkg/infrastructure/router"
)

type HomeController struct{}

func NewHomeController(r router.Router) {
	controller := &HomeController{}
	r.Get("/", controller.Index)
}

func (hc *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	page := renderer.GetPageTemplate("home")
	data := map[string]string{
		"Title": "Home",
		"Path":  "/",
	}
	page.ExecuteTemplate(w, "base", data)
}
