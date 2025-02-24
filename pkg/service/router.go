package service

import (
	"context"
	"net/http"
	"strings"

	"go-starter-template/pkg/middlewares"
)

type Router interface {
	http.Handler
	HandleFunc(string, http.HandlerFunc)
	Handle(string, http.Handler)
	Use(middlewares ...func(http.Handler) http.Handler)
	Wants(*http.Request, string) bool
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Route(pattern string, fn func(r Router)) Router
	// TODO: implement if grouping of endpoints is needed
	// refer to chi's doc for more interfaces. https://github.com/go-chi/chi
}

type Route struct {
	handler    http.Handler
	paramNames []string
	hasParams  bool
}

type NetServerMux struct {
	conf     *Config
	handlers map[string]*Route
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
		handlers: make(map[string]*Route),
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

func extractPathParams(pattern string) (cleanPattern string, paramNames []string) {
	parts := strings.Split(pattern, "/")
	cleanParts := make([]string, len(parts))
	paramNames = make([]string, 0)

	for i, part := range parts {
		if len(part) > 0 && part[0] == '{' && part[len(part)-1] == '}' {
			paramName := part[1 : len(part)-1]
			paramNames = append(paramNames, paramName)
			cleanParts[i] = "*"
		} else {
			cleanParts[i] = part
		}
	}

	return strings.Join(cleanParts, "/"), paramNames
}

func matchRoute(pattern string, paramNames []string, path string) (bool, map[string]string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	paramIndex := 0

	for i := 0; i < len(patternParts); i++ {
		if patternParts[i] == "*" {
			if paramIndex < len(paramNames) {
				params[paramNames[paramIndex]] = pathParts[i]
				paramIndex++
			}
		} else if patternParts[i] != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

func (n *NetServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := n.applyMiddlewares(n.mux)

	var route *Route
	key := r.URL.Path
	if r.Method != "GET" {
		key = r.Method + " " + r.URL.Path
	}

	route, exists := n.handlers[key]
	if !exists {
		for k, v := range n.handlers {
			if !v.hasParams {
				continue
			}

			methodMatches := true
			path := k
			if r.Method != "GET" {
				if !strings.HasPrefix(k, r.Method+" ") {
					continue
				}
				path = strings.TrimPrefix(k, r.Method+" ")
			}

			matched, params := matchRoute(path, v.paramNames, r.URL.Path)
			if matched && methodMatches {
				ctx := context.WithValue(r.Context(), "params", params)
				r = r.WithContext(ctx)
				route = v
				exists = true
				break
			}
		}
	}

	if exists {
		h := n.applyMiddlewares(route.handler)
		h.ServeHTTP(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func (n *NetServerMux) addToRoutes(method string, pattern string, handler http.HandlerFunc) {
	cleanPattern, paramNames := extractPathParams(pattern)
	hasParams := len(paramNames) > 0

	key := n.prefix + cleanPattern
	if method != "GET" {
		key = method + " " + n.prefix + cleanPattern
	}

	n.handlers[key] = &Route{
		handler:    handler,
		paramNames: paramNames,
		hasParams:  hasParams,
	}
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

func GetParam(r *http.Request, paramName string) string {
	ctx := r.Context()
	if params, ok := ctx.Value("params").(map[string]string); ok {
		return params[paramName]
	}
	return ""
}
