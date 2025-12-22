package expressish

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	router := New()
	if router == nil {
		t.Fatal("New() returned nil")
	}
	if router.routes == nil {
		t.Fatal("Router routes slice is nil")
	}
	if len(router.routes) != 0 {
		t.Fatalf("Expected empty routes slice, got %d routes", len(router.routes))
	}
}

func TestGET(t *testing.T) {
	router := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET response"))
	}

	router.GET("/test", handler)

	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}
	if router.routes[0].method != http.MethodGet {
		t.Errorf("Expected method %s, got %s", http.MethodGet, router.routes[0].method)
	}
	if router.routes[0].path != "/test" {
		t.Errorf("Expected path /test, got %s", router.routes[0].path)
	}
}

func TestPOST(t *testing.T) {
	router := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POST response"))
	}

	router.POST("/test", handler)

	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}
	if router.routes[0].method != http.MethodPost {
		t.Errorf("Expected method %s, got %s", http.MethodPost, router.routes[0].method)
	}
}

func TestPUT(t *testing.T) {
	router := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PUT response"))
	}

	router.PUT("/test", handler)

	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}
	if router.routes[0].method != http.MethodPut {
		t.Errorf("Expected method %s, got %s", http.MethodPut, router.routes[0].method)
	}
}

func TestDELETE(t *testing.T) {
	router := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DELETE response"))
	}

	router.DELETE("/test", handler)

	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}
	if router.routes[0].method != http.MethodDelete {
		t.Errorf("Expected method %s, got %s", http.MethodDelete, router.routes[0].method)
	}
}

func TestPATCH(t *testing.T) {
	router := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PATCH response"))
	}

	router.PATCH("/test", handler)

	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}
	if router.routes[0].method != http.MethodPatch {
		t.Errorf("Expected method %s, got %s", http.MethodPatch, router.routes[0].method)
	}
}

func TestServeHTTP_ExactPathMatch(t *testing.T) {
	router := New()

	router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	router.POST("/submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /hello matches",
			method:         http.MethodGet,
			path:           "/hello",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "POST /submit matches",
			method:         http.MethodPost,
			path:           "/submit",
			expectedStatus: http.StatusCreated,
			expectedBody:   "Created",
		},
		{
			name:           "GET /notfound does not match",
			method:         http.MethodGet,
			path:           "/notfound",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "POST /hello wrong method",
			method:         http.MethodPost,
			path:           "/hello",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "GET /hello/ with trailing slash does not match",
			method:         http.MethodGet,
			path:           "/hello/",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestServeHTTP_MultipleRoutes(t *testing.T) {
	router := New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	router.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("users list"))
	})

	router.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("create user"))
	})

	router.GET("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user profile"))
	})

	tests := []struct {
		method       string
		path         string
		expectedBody string
	}{
		{http.MethodGet, "/", "root"},
		{http.MethodGet, "/users", "users list"},
		{http.MethodPost, "/users", "create user"},
		{http.MethodGet, "/users/profile", "user profile"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestServeHTTP_ImplementsHTTPHandler(t *testing.T) {
	router := New()

	// Verify that Router implements http.Handler interface
	var _ http.Handler = router
}
