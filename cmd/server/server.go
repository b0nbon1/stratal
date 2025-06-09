package server

import (
	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/b0nbon1/stratal/api/routes"
)

type Server struct {
	store  *db.Queries
	router *gin.Engine
}

func NewServer(store *db.Queries) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	api := router.Group("/api")
	{
		routes.RegisterJobRoutes(api, store)
	}

	router.Static("/assets", "./client/public")
	router.NoRoute(func(c *gin.Context) {
		c.File("./client/dist/index.html")
	})

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
