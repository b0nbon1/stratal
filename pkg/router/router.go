// Package router provides a flexible HTTP router with support for parameterized routes,
// middleware chaining, and route grouping. It handles route conflicts intelligently by
// prioritizing static routes over parameterized ones.
//
// Example usage:
//
//	r := router.NewRouter()
//
//	// Static routes have higher priority
//	r.Get("/jobs", listJobsHandler)
//
//	// Parameterized routes have lower priority and can extract parameters
//	r.Get("/jobs/:id", getJobHandler)
//	r.Put("/jobs/:id", updateJobHandler)
//	r.Delete("/jobs/:id", deleteJobHandler)
//
//	// In your handler, extract the parameter:
//	func getJobHandler(w http.ResponseWriter, r *http.Request) {
//		id := router.GetParam(r, "id")
//		// use the id...
//	}
//
//	// Route groups work with parameterized routes too:
//	api := r.Group("/api/v1")
//	api.Get("/users/:userId/jobs/:jobId", getUserJobHandler)
package router

import (
	"context"
	"net/http"
	"sort"
	"strings"
)

// RouteHandler ensures handler has standard signature
type RouteHandler func(http.ResponseWriter, *http.Request)

// Middleware defines a standard middleware signature
type Middleware func(http.Handler) http.Handler

// Route represents a registered route with its pattern and metadata
type Route struct {
	Method     string
	Pattern    string
	Handler    http.Handler
	ParamNames []string
	IsParam    []bool
	Priority   int // Lower number = higher priority
}

// RouteParams holds extracted path parameters
type RouteParams map[string]string

// Context key for route parameters
type contextKey string

const ParamsKey contextKey = "route_params"

type Router struct {
	routes      *[]Route // Use pointer to share routes between grouped routers
	middlewares []Middleware
	prefix      string
}

// NewRouter returns a new top-level Router
func NewRouter() *Router {
	routes := make([]Route, 0)
	return &Router{
		routes: &routes, // Store pointer to routes slice
	}
}

// ServeHTTP makes Router satisfy http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle preflight requests (CORS OPTIONS)
    if req.Method == http.MethodOptions {
        // Run through middleware chain (so CORS headers get added)
        h := chainMiddlewares(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNoContent) // 204
        }), r.middlewares...)

        h.ServeHTTP(w, req)
        return
    }

	// Try to find a matching route
	if route, params := r.findMatchingRoute(req.Method, req.URL.Path); route != nil {
		// Add parameters to request context if any
		if len(params) > 0 {
			ctx := context.WithValue(req.Context(), ParamsKey, params)
			req = req.WithContext(ctx)
		}
		route.Handler.ServeHTTP(w, req)
		return
	}

	// No route found, return 404
	http.NotFound(w, req)
}

// GetParam extracts a path parameter from the request context
func GetParam(r *http.Request, key string) string {
	if params, ok := r.Context().Value(ParamsKey).(RouteParams); ok {
		return params[key]
	}
	return ""
}

// Group creates a subgroup with path prefix and optional middleware
func (r *Router) Group(prefix string, mws ...Middleware) *Router {
	return &Router{
		routes:      r.routes, // Share pointer to routes slice for parameterized routes
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

func (r *Router) Options(path string, mws ...interface{}) {
	r.handle("OPTIONS", path, mws...)
}

func (r *Router) Use(mws ...Middleware) {
	r.middlewares = append(r.middlewares, mws...)
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

	// Register all routes in our custom system for proper priority handling
	r.registerRoute(method, finalPath, finalHandler)
}

// registerRoute registers any route (static or parameterized) in our route system
func (r *Router) registerRoute(method, pattern string, handler http.Handler) {
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	paramNames := make([]string, len(segments))
	isParam := make([]bool, len(segments))
	priority := 0

	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			paramNames[i] = segment[1:] // Remove the ':' prefix
			isParam[i] = true
			priority += 10 // Parameterized segments have lower priority
		} else {
			isParam[i] = false
			priority += 1 // Static segments have higher priority
		}
	}

	route := Route{
		Method:     method,
		Pattern:    pattern,
		Handler:    handler,
		ParamNames: paramNames,
		IsParam:    isParam,
		Priority:   priority,
	}

	*r.routes = append(*r.routes, route)

	// Sort routes by priority (lower number = higher priority)
	sort.Slice(*r.routes, func(i, j int) bool {
		return (*r.routes)[i].Priority < (*r.routes)[j].Priority
	})
}

// findMatchingRoute finds the best matching route for the given method and path
func (r *Router) findMatchingRoute(method, path string) (*Route, RouteParams) {
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")

	for _, route := range *r.routes {
		if route.Method != method {
			continue
		}

		routeSegments := strings.Split(strings.Trim(route.Pattern, "/"), "/")

		// Check if segment count matches
		if len(pathSegments) != len(routeSegments) {
			continue
		}

		params := make(RouteParams)
		matches := true

		for i, routeSegment := range routeSegments {
			if route.IsParam[i] {
				// This is a parameter segment, extract the value
				params[route.ParamNames[i]] = pathSegments[i]
			} else {
				// This is a static segment, must match exactly
				if routeSegment != pathSegments[i] {
					matches = false
					break
				}
			}
		}

		if matches {
			return &route, params
		}
	}

	return nil, nil
}

// utility
func chainMiddlewares(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
