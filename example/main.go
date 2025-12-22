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

	// Register routes using Express-style API
	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Welcome to express-ish!\n")
	})

	router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello, World!\n")
	})

	router.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"users": ["alice", "bob", "charlie"]}`)
	})

	router.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"message": "User created successfully"}`)
	})

	router.PUT("/users/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "User 123 updated successfully"}`)
	})

	router.DELETE("/users/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "User 123 deleted successfully"}`)
	})

	router.PATCH("/users/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "User 123 partially updated"}`)
	})

	// Start the server
	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	log.Printf("Try visiting: http://localhost%s/\n", addr)
	log.Printf("              http://localhost%s/hello\n", addr)
	log.Printf("              http://localhost%s/users\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
