package expressish

import (
	"fmt"
	"net/http"
	"strings"
)

func New(opts ...Options) *App {
	app := &App{}

	mux := http.NewServeMux()
	app.mux = mux

	if len(opts) > 0 {
		app.opts = opts[0]
	}

	return app
}

func (a *App) Listen() error {
	addr := a.opts.Addr
	a.mux.Handle("/", http.HandlerFunc(a.serve))
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("listening on", addr)
	return http.ListenAndServe(addr, a.mux)
}

type paramsKey struct{}

func (a *App) serve(w http.ResponseWriter, r *http.Request) {
	reqParts := strings.Split(clearPath(r.URL.Path), "/")[1:]

	for _, route := range a.routes {
		if route.method != r.Method {
			continue
		}

		if len(reqParts) != len(route.parts) {
			continue
		}

		params := map[string]string{}
		matched := true

		for i, part := range route.parts {
			if strings.HasPrefix(part, ":") {
				params[part[1:]] = reqParts[i]
				continue
			}
			if part != reqParts[i] {
				matched = false
				break
			}
		}

		if matched {
			ctx := &Ctx{
				w:      w,
				r:      r,
				params: params,
			}
			route.handler(ctx)
			return
		}
	}

	http.NotFound(w, r)
}

func clearPath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	return strings.TrimSuffix(p, "/")
}

func (a *App) add(method, path string, fns ...any) {
	if len(fns) == 0 {
		panic("expressish: missing handler")
	}

	var handler Handler
	var mws []Middleware
	for i, fn := range fns {
		switch v := fn.(type) {

		case Middleware:
			mws = append(mws, v)

		case func(*Ctx, Next):
			mws = append(mws, Middleware(v))

		case Handler:
			if i != len(fns)-1 {
				panic("expressish: handler must be last")
			}
			handler = v

		case func(*Ctx):
			if i != len(fns)-1 {
				panic("expressish: handler must be last")
			}
			handler = Handler(v)

		default:
			panic(fmt.Sprintf("expressish: invalid argument %T", fn))
		}
	}

	if handler == nil {
		panic("expressish: no final handler")
	}

	clean := clearPath(path)
	parts := strings.Split(clean, "/")[1:]

	final := handler

	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		next := final
		final = func(c *Ctx) {
			mw(c, func() { next(c) })
		}
	}

	for i := len(a.mws) - 1; i >= 0; i-- {
		mw := a.mws[i]
		next := final
		final = func(c *Ctx) {
			mw(c, func() { next(c) })
		}
	}

	a.routes = append(a.routes, Route{
		method:  method,
		path:    clean,
		parts:   parts,
		handler: final,
	})
}

func (a *App) Get(path string, fns ...any) {
	a.add(http.MethodGet, path, fns...)
}

func (a *App) Post(path string, fns ...any) {
	a.add(http.MethodPost, path, fns...)
}

func (a *App) Delete(path string, fns ...any) {
	a.add(http.MethodDelete, path, fns...)
}

func Param(r *http.Request, key string) (string, bool) {
	params, ok := r.Context().Value(paramsKey{}).(map[string]string)
	if !ok {
		return "", false
	}
	v, ok := params[key]
	return v, ok
}

func (a *App) Use(mw Middleware) {
	a.mws = append(a.mws, mw)
}
