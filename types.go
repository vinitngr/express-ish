package expressish

import (
	"net/http"
	"time"
)

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
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type App struct {
	opts   Options
	routes []Route
	mws    []Middleware
	mux    *http.ServeMux
}