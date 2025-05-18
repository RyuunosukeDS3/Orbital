package app

import (
	"github.com/gin-gonic/gin"
)

// Server struct represents the Gin application server
type Server struct {
	router *gin.Engine
}

// NewServer creates a new Gin server
func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

// Run starts the Gin server
func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
