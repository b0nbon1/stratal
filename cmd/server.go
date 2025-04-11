package cmd

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	// store  db.Store
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{
		// store:  store,
	}
	router := gin.Default()

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}