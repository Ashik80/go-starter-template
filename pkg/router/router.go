package router

import (
	"context"
	"net/http"
	"strings"
)

type Router interface {
	http.Handler
	HandleFunc(string, http.HandlerFunc)
	Handle(string, http.Handler)
	Use(middlewares ...func(http.Handler) http.Handler)
	With(middlewares ...func(http.Handler) http.Handler) Router
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

type contextKey string

const paramsContextKey contextKey = "params"

type Route struct {
	handler    http.Handler
	paramNames []string
	hasParams  bool
}

type NetServerMux struct {
	handlers map[string]*Route
	prefix   string
	mux      *http.ServeMux
	mws      []func(http.Handler) http.Handler
}

func NewNetServerMux() *NetServerMux {
	mux := http.NewServeMux()
	n := &NetServerMux{
		prefix:   "",
		mux:      mux,
		mws:      make([]func(http.Handler) http.Handler, 0),
		handlers: make(map[string]*Route),
	}

	return n
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

			path := k
			if r.Method != "GET" {
				if !strings.HasPrefix(k, r.Method+" ") {
					continue
				}
				path = strings.TrimPrefix(k, r.Method+" ")
			}

			matched, params := matchRoute(path, v.paramNames, r.URL.Path)
			if matched {
				ctx := context.WithValue(r.Context(), paramsContextKey, params)
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
	n.mws = append(n.mws, middlewares...)
}

func (n *NetServerMux) With(middlewares ...func(http.Handler) http.Handler) Router {
	child := &NetServerMux{
		prefix:   n.prefix,
		mux:      n.mux,
		handlers: n.handlers,
		mws:      append(n.mws, middlewares...),
	}
	return child
}

func (n *NetServerMux) Route(pattern string, fn func(r Router)) Router {
	subMux := NewNetServerMux()
	subMux.prefix = n.prefix + pattern
	subMux.mws = append(subMux.mws, n.mws...)

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
	if params, ok := ctx.Value(paramsContextKey).(map[string]string); ok {
		return params[paramName]
	}
	return ""
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

	// for i := 0; i < len(patternParts); i++ {
	for i := range patternParts {
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
