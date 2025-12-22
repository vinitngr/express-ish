package expressish

import "net/http"

type Handler func(*Ctx)

type Route struct {
	method  string
	path    string
	handler Handler
	parts []string
}

type Options struct {
	Addr string
}

type App struct {
	opts   Options
	routes []Route
	mux *http.ServeMux
}