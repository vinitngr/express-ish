package main

import expressish "github.com/vinitngr/express-ish"

func main() {
	app := expressish.New(expressish.Options{
		Addr: ":8080",
	})

	app.Get("/health", func(c *expressish.Ctx) {
		c.Text(200, "ok")
	})

	app.Get("/users/:id", func(c *expressish.Ctx) {
		id := c.Param("id")
		c.Text(200, "user "+id)
	})

	app.Listen()
}
