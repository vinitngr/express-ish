package expressish

import "net/http"

type Handler func(*Ctx)

type Next func()
type Middleware func(*Ctx, Next)

type Route struct {
	method  string
	path    string
	handler Handler
	parts   []string
	mws     []Middleware
}

type Options struct {
	Addr string
}

type App struct {
	opts   Options
	routes []Route
	mws    []Middleware
	mux    *http.ServeMux
}