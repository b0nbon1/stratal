package api

import (
	"encoding/json"
	"net/http"

	"github.com/b0nbon1/stratal/internal/api/middleware"
	"github.com/b0nbon1/stratal/pkg/router"
)

func (hs *HTTPServer) registerRoutes() *router.Router {
	r := router.NewRouter()
	r.Use(middleware.CORSMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, 200, map[string]interface{}{
			"status":  "healthy",
			"service": "stratal",
		})
	})

	api := r.Group("/api", middleware.Logging())

	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, 200, map[string]interface{}{
			"message": "API is running",
		})
	})

	v1 := api.Group("/v1")

	v1.Post("/jobs", hs.CreateJob)
	v1.Get("/jobs", hs.ListJobs)
	v1.Get("/jobs/:id", hs.GetJob)

	v1.Post("/job-runs", hs.CreateJobRun)
	v1.Get("/job-runs", hs.GetJobRun)
	v1.Get("/job-runs/:id", hs.GetJobRun)
	
	// Job run control endpoints
	v1.Post("/job-runs/:id/pause", hs.PauseJobRun)
	v1.Post("/job-runs/:id/resume", hs.ResumeJobRun)
	v1.Get("/job-runs/paused", hs.GetPausedJobRuns)

	v1.Post("/secrets", hs.CreateSecret)
	v1.Get("/secrets", hs.ListSecrets)

	// Add log routes
	if hs.logSystem != nil {
		logsHandler := NewLogsHandler(hs.store, hs.logSystem.GetStreamer())
		logsHandler.RegisterLogRoutes(v1)
	}

	return r
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
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

func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
