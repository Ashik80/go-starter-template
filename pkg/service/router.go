package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	HandleFunc(string, http.HandlerFunc)
	Handle(string, http.Handler)
	WithPathParams(*http.Request) map[string]string
}

type ChiServerMux struct {
	router *chi.Mux
}

func NewChiServerMux() *ChiServerMux {
	return &ChiServerMux{router: chi.NewRouter()}
}

func (c *ChiServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

func (c *ChiServerMux) HandleFunc(path string, handler http.HandlerFunc) {
	c.router.HandleFunc(path, handler)
}

func (c *ChiServerMux) Handle(path string, handler http.Handler) {
	c.router.Handle(path, handler)
}

func (c *ChiServerMux) WithPathParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	ctx := chi.RouteContext(r.Context())
	if ctx != nil {
		for i, key := range ctx.URLParams.Keys {
			params[key] = ctx.URLParams.Values[i]
		}
	}
	return params
}
