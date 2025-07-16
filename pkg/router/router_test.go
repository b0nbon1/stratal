package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteConflictResolution(t *testing.T) {
	r := NewRouter()

	// Add a static route first
	r.Get("/jobs", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("list jobs"))
	})

	// Add a parameterized route that could conflict
	r.Get("/jobs/:id", func(w http.ResponseWriter, req *http.Request) {
		id := GetParam(req, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("job id: " + id))
	})

	tests := []struct {
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"/jobs", 200, "list jobs"},       // Should match static route
		{"/jobs/123", 200, "job id: 123"}, // Should match parameterized route
		{"/jobs/abc", 200, "job id: abc"}, // Should match parameterized route
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.path, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != test.expectedStatus {
			t.Errorf("Expected status %d for path %s, got %d", test.expectedStatus, test.path, w.Code)
		}

		if w.Body.String() != test.expectedBody {
			t.Errorf("Expected body '%s' for path %s, got '%s'", test.expectedBody, test.path, w.Body.String())
		}
	}
}

func TestParameterExtraction(t *testing.T) {
	r := NewRouter()

	r.Get("/users/:userId/jobs/:jobId", func(w http.ResponseWriter, req *http.Request) {
		userId := GetParam(req, "userId")
		jobId := GetParam(req, "jobId")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user: " + userId + ", job: " + jobId))
	})

	req := httptest.NewRequest("GET", "/users/42/jobs/abc123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := "user: 42, job: abc123"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}

func TestRouteGroups(t *testing.T) {
	r := NewRouter()

	api := r.Group("/api/v1")
	api.Get("/jobs/:id", func(w http.ResponseWriter, req *http.Request) {
		id := GetParam(req, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("api job: " + id))
	})

	req := httptest.NewRequest("GET", "/api/v1/jobs/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := "api job: test"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}

func TestRoutePriority(t *testing.T) {
	r := NewRouter()

	// Add parameterized route first
	r.Get("/special/:type", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("parameterized"))
	})

	// Add static route after - should have higher priority
	r.Get("/special/admin", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("static"))
	})

	// Test that static route takes precedence
	req := httptest.NewRequest("GET", "/special/admin", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Should match the static route since it has higher priority
	expected := "static"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s' (static route should have priority), got '%s'", expected, w.Body.String())
	}
}
