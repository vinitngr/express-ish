package main

import (
	"fmt"
	"os"

	expressish "github.com/vinitngr/express-ish"
	"github.com/vinitngr/express-ish/middleware"
)

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

	app.Use(middleware.Logger)

	app.Use(middleware.DotEnv(".example.env"))
	fmt.Println(os.Getenv("SOMETHING"))

	app.Use(middleware.CORS([]string{
		"http://localhost:3000",
	}))
	
	app.Use(middleware.Recovery)

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

	app.Get("/crash", func(c *expressish.Ctx) {
		panic("intentional crash")
	})

	app.Post(
		"/users/:id/action",
		Auth,
		func(c *expressish.Ctx) {

			id := c.Param("id")

			debug := c.Query("debug")
			mode := c.Query("mode")

			if c.ReqHeader("Content-Type") != "application/json" {
				c.Text(415, "content-type must be application/json")
				return
			}

			body := c.Body()
			action, ok := body["action"].(string)
			if !ok || action == "" {
				c.Text(400, "missing action")
				return
			}

			c.Set("X-User-ID", id)
			if debug == "1" {
				c.Set("X-Debug-Mode", "true")
			}

			if action == "redirect" {
				c.Redirect("/health", 302)
				return
			}

			c.JSON(200, map[string]any{
				"id":     id,
				"mode":   mode,
				"debug":  debug == "1",
				"action": action,
			})
		},
	)

	app.Listen()
}
