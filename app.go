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
	a.mux.Handle("/" , http.HandlerFunc(a.serve))
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("listening on" , addr)
	return http.ListenAndServe(addr , a.mux)
}

type paramsKey struct{}
func (a *App) serve(w http.ResponseWriter, r *http.Request){
	reqParts := strings.Split(clearPath(r.URL.Path), "/")[1:]

	for _ , route := range a.routes {
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
				W:      w,
				R:      r,
				params: params,
			}
			route.handler(ctx)
			return
		}

		http.NotFound(w , r)
	}
}

func clearPath(p string) string {
	if p == "" { return "/" }
	if !strings.HasPrefix(p , "/"){ p = "/" + p }

	return strings.TrimSuffix(p , "/")
}

func (a *App) add(method, path string, h Handler) {
	clean := clearPath(path)
	parts := strings.Split(clean, "/")[1:]

	a.routes = append(a.routes, Route{
		method:  method,
		path:    clean,
		parts:   parts,
		handler: h,
	})
}

func (a *App) Get(path string, h Handler) {
	a.add(http.MethodGet, path, h)
}

func (a *App) Post(path string, h Handler) {
	a.add(http.MethodPost, path, h)
}

func (a *App) Delete(path string, h Handler) {
	a.add(http.MethodDelete, path, h)
}

func Param(r *http.Request, key string) (string, bool) {
    params, ok := r.Context().Value(paramsKey{}).(map[string]string)
    if !ok {
        return "", false
    }
    v, ok := params[key]
    return v, ok
}