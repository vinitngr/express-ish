# express-ish

`express-ish is a lightweight HTTP routing and middleware library for Go, inspired by Express.js and built on top of the standard net/http package.`

It focuses on **developer ergonomics** (middleware, params, body parsing, helpers) while keeping Go’s execution model explicit and predictable.

This project is learning-driven but fully usable for real services.

---

## Features

* Express-style routing (`app.Get`, `app.Post`)
* Global and route-level middleware
* Custom request context (`Ctx`)
* Path params (`/users/:id`)
* Query parsing
* JSON and form body parsing
* JSON, text, redirect responses
* Request-scoped data sharing (`Set` / `Get`)
* Timeout, recovery, CORS, dotenv middleware
* Built entirely on Go standard library

---

## Installation

```bash
go get github.com/vinitngr/express-ish
```

---

## Basic Usage

```go
app := expressish.New(expressish.Options{
	Addr: ":8080",
})

app.Get("/health", func(c *expressish.Ctx) {
	c.Text(200, "ok")
})

app.Listen()
```

---

## Middleware

Middleware follows an Express-style flow using `next()`.

```go
func Auth(c *expressish.Ctx, next expressish.Next) {
	if c.Param("id") == "1" {
		c.Text(401, "unauthorized")
		return
	}
	next()
}
```

Register globally or per route:

```go
app.Use(middleware.Logger)

app.Get("/users/:id", Auth, handler)
```

Execution order:

```
global middleware → route middleware → handler
```

---

## Example (Realistic)

A trimmed example demonstrating routing, middleware, params, query, body parsing, and responses:

```go
app.Use(middleware.Logger)
app.Use(middleware.Recovery)

app.Get("/users/:id", Auth, func(c *expressish.Ctx) {
	id := c.Param("id")
	debug := c.Query("debug")

	c.Set("X-User-ID", id)

	c.JSON(200, map[string]any{
		"id":    id,
		"debug": debug == "1",
	})
})

app.Post("/users/:id/action", Auth, func(c *expressish.Ctx) {
	if c.ReqHeader("Content-Type") != "application/json" {
		c.Text(415, "content-type must be application/json")
		return
	}

	type Req struct {
		Action string `json:"action" form:"action"`
	}

	var body Req
	if err := c.Body(&body); err != nil {
		c.Text(400, err.Error())
		return
	}

	if body.Action == "redirect" {
		c.Redirect("/health", 302)
		return
	}

	c.JSON(200, map[string]any{
		"id":     c.Param("id"),
		"action": body.Action,
	})
})
```

---

## Request Helpers

### Params

```go
c.Param("id")
```

### Query

```go
c.Query("page")
c.QueryInt("page", 1)
c.QueryBool("debug", false)
```

### Body (JSON / form, lazy)

```go
var data MyStruct
err := c.Body(&data)
```

Parsed only when accessed.

---

## Response Helpers

```go
c.Text(200, "ok")
c.JSON(200, data)
c.Redirect("/health", 302)
c.Status(204)
```

---

## Context Data Sharing

Middleware and handlers share request-scoped data:

```go
c.Set("user", "admin")
user := c.Get("user")
```

---

## Design Notes

* Uses Go stdlib internally (`net/http`, `net/url`, `encoding/json`)
* No hidden goroutines or background work
* Lazy parsing (query/body parsed only when used)
* Express-style middleware implemented via function composition
* `next()` is framework-level, not `net/http` based

---

## Status

Experimental but functional.

APIs may evolve as features are added incrementally.

---

## License

MIT
