package server

import (
	"github.com/DevHeaven/db/internal/middleware"
	"github.com/DevHeaven/db/internal/router"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	port   string
	router router.Router
}

func NewServer(port string, engine *gin.Engine, routes router.Router) Server {
	return Server{
		Engine: engine,
		port:   port,
		router: routes,
	}
}

func (s *Server) Run() error {
	s.Engine.Use(middleware.CORSmanager)

	s.Engine.GET("/health", func(ctx *gin.Context) {
		s.router.Health(ctx)
	})

	s.Engine.POST("/pep1", func(ctx *gin.Context) {
		s.router.Pep1(ctx)
	})

	return s.Engine.Run(s.port)
}
