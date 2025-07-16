package api

import (
	"encoding/json"
	"net/http"

	"github.com/b0nbon1/stratal/internal/api/middleware"
	"github.com/b0nbon1/stratal/pkg/router"
)

func (hs *HTTPServer) registerRoutes() *router.Router {

	r := router.NewRouter()

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, 200, map[string]interface{}{
			"status":  "healthy",
			"service": "stratal",
		})
	})

	// API routes with logging middleware
	api := r.Group("/api", middleware.Logging())

	// API health check
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, 200, map[string]interface{}{
			"message": "API is running",
		})
	})

	// V1 API routes
	v1 := api.Group("/v1")

	// Job routes
	v1.Post("/jobs", hs.CreateJob)
	v1.Get("/jobs", hs.ListJobs)
	v1.Get("/jobs/:id", hs.GetJob)

	// Job run routes
	v1.Post("/job-runs", hs.CreateJobRun)
	v1.Get("/job-runs", hs.GetJobRun)
	v1.Get("/job-runs/:id", hs.GetJobRun)

	return r
}

func parseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string, details ...interface{}) {
	response := map[string]interface{}{
		"error": message,
	}

	if len(details) > 0 {
		response["details"] = details[0]
	}

	respondJSON(w, status, response)
}
