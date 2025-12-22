package main

import (
	expressish "github.com/vinitngr/express-ish"
	"fmt"
)

func Logger(c *expressish.Ctx, next expressish.Next) {
	fmt.Println("logger: before")
	next()
	fmt.Println("logger: after")
}

func Auth(c *expressish.Ctx, next expressish.Next) {
	user := c.Param("id")
	
	if user == "1" {
		c.Text(401, "unauthorized")
		return
	}
	next()
}
func main() {
	app := expressish.New(expressish.Options{
		Addr: ":8080",
	})

	app.Use(Logger)
	app.Get("/health", func(c *expressish.Ctx) {
		c.Text(200, "ok")
	})

	app.Get(
		"/users/:id",
		Auth,
		func(c *expressish.Ctx) {
			id := c.Param("id")
			c.Text(200, "user "+id)
		},
	)

	app.Listen()
}
