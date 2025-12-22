package expressish

import "net/http"

type Handler func(http.ResponseWriter, *http.Request)

type Route struct {
	method  string
	path    string
	handler Handler
}

type Options struct {
	Addr string
}

type App struct {
	opts   Options
	routes []Route
	mux *http.ServeMux
}