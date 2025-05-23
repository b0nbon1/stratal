package api

import (
	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Queries
	router *gin.Engine
}

func NewServer(store db.Queries) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/jobs", server.createJobRequest)
	// router.GET("/jobs/:id", server.getJobRequest)
	router.GET("/jobs", server.listJobsRequest)
	// router.PUT("/jobs/:id", server.updateJobRequest)
	// router.DELETE("/jobs/:id", server.deleteJobRequest)
	// router.POST("/jobs/:id/run", server.runJobRequest)
	// router.POST("/jobs/:id/stop", server.stopJobRequest)
	// router.POST("/jobs/:id/retry", server.retryJobRequest)
	// router.POST("/jobs/:id/force", server.forceJobRequest)
	// router.POST("/jobs/:id/force-stop", server.forceStopJobRequest)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
