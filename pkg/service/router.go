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
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)

	Route(pattern string, fn func(r Router)) Router
}

type NetServerMux struct {
	conf     *Config
	handlers map[string]http.Handler
	prefix   string
	mux      *http.ServeMux
	mws      []middlewares.MiddlewareFunc
}

func NewNetServerMux(conf *Config) *NetServerMux {
	mux := http.NewServeMux()
	n := &NetServerMux{
		conf:     conf,
		prefix:   "",
		mux:      mux,
		mws:      []middlewares.MiddlewareFunc{},
		handlers: make(map[string]http.Handler),
	}

	n.Use(middlewares.Logger)
	csrfMiddleware := middlewares.CSRFMiddleware(conf.CSRFAuthKey)
	n.Use(csrfMiddleware)
	n.Use(middlewares.EnableCors(conf.AllowedOrigins))

	return n
}

func (n *NetServerMux) applyMiddlewares(handler http.Handler) http.Handler {
	for _, mw := range n.mws {
		handler = mw(handler)
	}
	return handler
}

func (n *NetServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := n.applyMiddlewares(n.mux)

	var route http.Handler
	if r.Method == "GET" {
		route = n.handlers[r.URL.Path]
	} else {
		route = n.handlers[r.Method+" "+r.URL.Path]
	}

	if route != nil {
		route = n.applyMiddlewares(route)
		route.ServeHTTP(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func (n *NetServerMux) addToRoutes(method string, pattern string, handler http.HandlerFunc) {
	if method == "GET" {
		n.handlers[n.prefix+pattern] = handler
		return
	}
	n.handlers[method+" "+n.prefix+pattern] = handler
}

func (n *NetServerMux) Get(pattern string, handler http.HandlerFunc) {
	method := "GET"
	n.addToRoutes(method, pattern, handler)
}

func (n *NetServerMux) Post(pattern string, handler http.HandlerFunc) {
	method := "POST"
	n.addToRoutes(method, pattern, handler)
}

func (n *NetServerMux) Put(pattern string, handler http.HandlerFunc) {
	method := "PUT"
	n.addToRoutes(method, pattern, handler)
}

func (n *NetServerMux) Patch(pattern string, handler http.HandlerFunc) {
	method := "PATCH"
	n.addToRoutes(method, pattern, handler)
}

func (n *NetServerMux) Delete(pattern string, handler http.HandlerFunc) {
	method := "DELETE"
	n.addToRoutes(method, pattern, handler)
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

func (n *NetServerMux) Route(pattern string, fn func(r Router)) Router {
	subMux := NewNetServerMux(n.conf)
	subMux.prefix = n.prefix + pattern
	fn(subMux)
	for k, v := range subMux.handlers {
		if k[len(k)-1] == '/' {
			k = k[:len(k)-1]
		}
		n.handlers[k] = v
	}
	return subMux
}
