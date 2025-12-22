package main

import (
	"fmt"
	"net/http"
	"github.com/vinitngr/express-ish"
)

func main() {
	app := expressish.New(expressish.Options{
		Addr: ":8080",
	})

	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("health handler hit")
		w.Write([]byte("ok"))
	})

	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("root handler hit")
		w.Write([]byte("hello"))
	})

	app.Listen()
}
