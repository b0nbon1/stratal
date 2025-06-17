package api

import (
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/internal/api/middleware"
	"github.com/b0nbon1/stratal/internal/api/router"
)

func (httpServer *HTTPServer) registerRoutes() *router.Router {
	r := router.NewRouter()

	api := r.Group("/api", middleware.Logging())

	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("We're here almost done, guess healthcheck is okay")
		respondJSON(w, 200, map[string]interface{}{
			"message": "OK",
		})
	})

	return r
}

