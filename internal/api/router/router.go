package router

import (
	"net/http"
	"strings"
)

// RouteHandler ensures handler has standard signature
type RouteHandler func(http.ResponseWriter, *http.Request)

// Middleware defines a standard middleware signature
type Middleware func(http.Handler) http.Handler

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
	prefix      string
}

// NewRouter returns a new top-level Router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// ServeHTTP makes Router satisfy http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Group creates a subgroup with path prefix and optional middleware
func (r *Router) Group(prefix string, mws ...Middleware) *Router {
	return &Router{
		mux:         r.mux,
		prefix:      r.prefix + strings.TrimSuffix(prefix, "/"),
		middlewares: append(r.middlewares, mws...),
	}
}

func (r *Router) Get(path string, mws ...interface{}) {
	r.handle("GET", path, mws...)
}

func (r *Router) Post(path string, mws ...interface{}) {
	r.handle("POST", path, mws...)
}

func (r *Router) Patch(path string, mws ...interface{}) {
	r.handle("PATCH", path, mws...)
}

func (r *Router) Put(path string, mws ...interface{}) {
	r.handle("PUT", path, mws...)
}

func (r *Router) Delete(path string, mws ...interface{}) {
	r.handle("DELETE", path, mws...)
}

func (r *Router) handle(method, path string, args ...interface{}) {
	finalPath := r.prefix + path

	if len(args) == 0 {
		panic("handler is required")
	}

	// extract the handler
	last := args[len(args)-1]

	var handler http.Handler

	switch h := last.(type) {
	case RouteHandler:
		handler = http.HandlerFunc(h)
	case http.HandlerFunc:
		handler = h
	case func(http.ResponseWriter, *http.Request):
		handler = http.HandlerFunc(h)
	default:
		panic("last argument must be a valid HTTP handler")
	}

	// collect middlewares
	var middlewares []Middleware
	for _, arg := range args[:len(args)-1] {
		if mw, ok := arg.(Middleware); ok {
			middlewares = append(middlewares, mw)
		} else {
			panic("middleware must be of type Middleware")
		}
	}

	finalHandler := chainMiddlewares(handler, append(r.middlewares, middlewares...)...)

	r.mux.Handle(finalPath, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		finalHandler.ServeHTTP(w, req)
	}))
}


// utility
func chainMiddlewares(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
