package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/internal/api/middleware"
	"github.com/b0nbon1/stratal/internal/api/router"
)

func (hs *HTTPServer) registerRoutes() *router.Router {
	r := router.NewRouter()

	api := r.Group("/api", middleware.Logging())

	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("We're here almost done, guess healthcheck is okay")
		respondJSON(w, 200, map[string]interface{}{
			"message": "OK",
		})
	})

	v1 := r.Group("/v1")
	v1.Post("/jobs", hs.CreateJob)

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

