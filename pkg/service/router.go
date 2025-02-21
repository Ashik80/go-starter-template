package service

import (
	"net/http"

	"go-starter-template/pkg/middlewares"
)

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	HandleFunc(string, http.HandlerFunc)
	Handle(string, http.Handler)
	Use(middlewares ...func(http.Handler) http.Handler)
	Wants(*http.Request, string) bool
	// TODO: implement if grouping of endpoints is needed
	// refer to chi's doc for more interfaces. https://github.com/go-chi/chi
	// Route(pattern string, fn func(r Router)) Router
}

type NetServerMux struct {
	mux *http.ServeMux
	mws []middlewares.MiddlewareFunc
}

func NewNetServerMux(conf *Config) *NetServerMux {
	mux := http.NewServeMux()
	n := &NetServerMux{
		mux: mux,
		mws: []middlewares.MiddlewareFunc{},
	}

	n.Use(middlewares.Logger)
	csrfMiddleware := middlewares.CSRFMiddleware(conf.CSRFAuthKey)
	n.Use(csrfMiddleware)
	n.Use(middlewares.EnableCors(conf.AllowedOrigins))

	return n
}

func (n *NetServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	firstMw := n.mws[0]
	h := firstMw(n.mux)

	for _, mw := range n.mws[1:] {
		h = mw(h)
	}

	h.ServeHTTP(w, r)
}

func (n *NetServerMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	n.mux.HandleFunc(pattern, handler)
}

func (n *NetServerMux) Handle(pattern string, handler http.Handler) {
	n.mux.Handle(pattern, handler)
}

func (n *NetServerMux) Wants(r *http.Request, accept string) bool {
	a := r.Header.Get("Accept")
	return a == accept
}

func (n *NetServerMux) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		n.mws = append(n.mws, mw)
	}
}
