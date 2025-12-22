# express-ish

A lightweight HTTP router for Go, built on top of Go's standard `net/http` package. Inspired by Express.js, it provides a simple and familiar API for building web applications while keeping Go's native behavior transparent.

## Features

- 🎯 **Single Entry Point**: Router implements `http.Handler` interface
- 💾 **In-Memory Route Storage**: Routes are stored in a simple slice for easy understanding
- 🔍 **Exact Path Matching**: Start with the fundamentals - exact path matching only
- 🚀 **Express-Style API**: Familiar method-based routing (GET, POST, PUT, DELETE, PATCH)
- 📦 **Zero Dependencies**: Only uses Go's standard library
- 🧠 **Learning-Focused**: Clear, readable code designed for understanding HTTP routing

## Installation

```bash
go get github.com/vinitngr/express-ish
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    "github.com/vinitngr/express-ish"
)

func main() {
    // Create a new router
    router := expressish.New()
    
    // Register routes
    router.GET("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to express-ish!")
    })
    
    router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    router.POST("/users", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintf(w, `{"message": "User created"}`)
    })
    
    // Start the server
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## API Reference

### Router

#### `New() *Router`

Creates and returns a new Router instance.

```go
router := expressish.New()
```

#### HTTP Method Functions

Register routes for different HTTP methods:

- `GET(path string, handler HandlerFunc)` - Register GET route
- `POST(path string, handler HandlerFunc)` - Register POST route
- `PUT(path string, handler HandlerFunc)` - Register PUT route
- `DELETE(path string, handler HandlerFunc)` - Register DELETE route
- `PATCH(path string, handler HandlerFunc)` - Register PATCH route

```go
router.GET("/users", getUsersHandler)
router.POST("/users", createUserHandler)
router.PUT("/users/123", updateUserHandler)
router.DELETE("/users/123", deleteUserHandler)
router.PATCH("/users/123", patchUserHandler)
```

### HandlerFunc

Handler functions have the same signature as `http.HandlerFunc`:

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

## How It Works

### Request Flow

1. **Route Registration**: Routes are stored in-memory as a slice of route structs
2. **Request Dispatch**: The router's `ServeHTTP` method implements the `http.Handler` interface
3. **Route Matching**: When a request arrives, routes are matched using:
   - HTTP method (GET, POST, etc.)
   - Exact path matching (no wildcards or parameters yet)
4. **Handler Execution**: If a match is found, the corresponding handler is called
5. **404 Response**: If no match is found, `http.NotFound` is called

### Example: Understanding the Flow

```go
router := expressish.New()

// Step 1: Register a route
router.GET("/api/users", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Users list"))
})
// Internal: route{method: "GET", path: "/api/users", handler: func} added to routes slice

// Step 2: Start the server
http.ListenAndServe(":8080", router)
// router is passed as http.Handler (it implements ServeHTTP)

// Step 3: Request arrives (GET /api/users)
// router.ServeHTTP is called by net/http

// Step 4: Route matching
// - Checks: method == "GET" && path == "/api/users"
// - Match found! Handler is executed

// Step 5: Handler writes response
// "Users list" is sent back to client
```

## Example Application

See the [example](./example) directory for a complete working example:

```bash
cd example
go run main.go
```

Then test the endpoints:

```bash
# GET requests
curl http://localhost:8080/
curl http://localhost:8080/hello
curl http://localhost:8080/users

# POST request
curl -X POST http://localhost:8080/users

# PUT request
curl -X PUT http://localhost:8080/users/123

# DELETE request
curl -X DELETE http://localhost:8080/users/123

# PATCH request
curl -X PATCH http://localhost:8080/users/123
```

## Exact Path Matching

Currently, express-ish uses **exact path matching** only:

- ✅ `/users` matches `/users`
- ❌ `/users` does NOT match `/users/`
- ❌ `/users` does NOT match `/users/123`
- ❌ No wildcard patterns
- ❌ No path parameters

This is intentional - the project starts simple and will gradually add features to help understand routing complexity step by step.

## Testing

Run the tests:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -cover
```

## Design Philosophy

express-ish is a **learning project** that prioritizes:

1. **Clarity over Performance**: Code is written to be readable and understandable
2. **Gradual Complexity**: Features are added step by step, starting from the basics
3. **Transparency**: Built on top of Go's `net/http` without hiding core behavior
4. **Simplicity**: Minimal abstractions, straightforward implementation

This is NOT meant to be a production router. For production use, consider:
- [gorilla/mux](https://github.com/gorilla/mux)
- [chi](https://github.com/go-chi/chi)
- [gin](https://github.com/gin-gonic/gin)

## Future Enhancements (Learning Path)

Potential features to add for learning:

- [ ] Path parameters (e.g., `/users/:id`)
- [ ] Query string parsing
- [ ] Middleware support
- [ ] Route groups
- [ ] Wildcard routes
- [ ] Static file serving
- [ ] Request context
- [ ] Pattern matching improvements

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

This is a learning project, but contributions that maintain the educational focus are welcome! Please ensure:

- Code is clear and well-commented
- Changes are minimal and focused
- Tests are included
- Documentation is updated

## Acknowledgments

Inspired by [Express.js](https://expressjs.com/) and built to understand Go's `net/http` package better.