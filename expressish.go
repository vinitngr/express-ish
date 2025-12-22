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

func (a *App) Get(path string , h Handler){
	a.routes = append(a.routes , Route{
		method : http.MethodGet,
		path : clearPath(path),
		handler : h,
	})
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

func (a *App) serve(w http.ResponseWriter, r *http.Request){
	reqPath := clearPath(r.URL.Path)

	fmt.Println("incoming" , r.Method , reqPath)

	for _ , route := range a.routes {
		if route.method != r.Method {
			continue 
		}

		if route.path == reqPath {
				fmt.Println("matched : " , route.path)
				route.handler(w , r)
				return
		}

		http.NotFound(w , r)
	}
}

func clearPath(p string) string {
	if p == "" {
		return "/"
	}

	if !strings.HasPrefix(p , "/"){
		p = "/" + p
	}

	return strings.TrimSuffix(p , "/")
}