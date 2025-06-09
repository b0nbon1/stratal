package routes

import (
	"github.com/b0nbon1/stratal/api/handlers"
	"github.com/b0nbon1/stratal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterJobRoutes(router *gin.Engine, store *db.Queries) {
	jobHandler := handlers.NewJobHandler(store)

	jobRoutes := router.Group("/jobs")
	{
		jobRoutes.POST("/", jobHandler.CreateJob)
		jobRoutes.GET("/:id", jobHandler.GetJobRequest)
		jobRoutes.GET("/", jobHandler.ListJobs)
	}
}