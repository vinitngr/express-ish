package expressish

import (
	"net/http"
)

// HandlerFunc defines the function signature for route handlers
type HandlerFunc func(http.ResponseWriter, *http.Request)

// route represents a single route with its HTTP method and handler
type route struct {
	method  string
	path    string
	handler HandlerFunc
}

// Router is the main router structure that stores routes in memory
type Router struct {
	routes []route
}

// New creates and returns a new Router instance
func New() *Router {
	return &Router{
		routes: make([]route, 0),
	}
}

// GET registers a GET route with the given path and handler
func (r *Router) GET(path string, handler HandlerFunc) {
	r.addRoute(http.MethodGet, path, handler)
}

// POST registers a POST route with the given path and handler
func (r *Router) POST(path string, handler HandlerFunc) {
	r.addRoute(http.MethodPost, path, handler)
}

// PUT registers a PUT route with the given path and handler
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.addRoute(http.MethodPut, path, handler)
}

// DELETE registers a DELETE route with the given path and handler
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.addRoute(http.MethodDelete, path, handler)
}

// PATCH registers a PATCH route with the given path and handler
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.addRoute(http.MethodPatch, path, handler)
}

// addRoute adds a route to the in-memory route storage
func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  method,
		path:    path,
		handler: handler,
	})
}

// ServeHTTP implements the http.Handler interface, making Router a single entry point
// It performs exact path matching to dispatch requests to registered handlers
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		// Exact path matching
		if route.method == req.Method && route.path == req.URL.Path {
			route.handler(w, req)
			return
		}
	}

	// No matching route found
	http.NotFound(w, req)
}
