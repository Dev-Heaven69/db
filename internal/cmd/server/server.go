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

	s.Engine.GET("/api/health", func(ctx *gin.Context) {
		s.router.Health(ctx)
	})

	s.Engine.POST("/api/scandb", func(ctx *gin.Context) {
		s.router.ScanDB(ctx)
	})

	s.Engine.POST("/api/getPersonalEmail", func(ctx *gin.Context) {
		s.router.GetPersonalEmail(ctx)
	})

	s.Engine.POST("/api/getProfessionalEmails", func(ctx *gin.Context) {
		s.router.GetProfessionalEmails(ctx)
	})

	s.Engine.POST("/api/getBothEmails", func(ctx *gin.Context) {
		s.router.GetBothEmails(ctx)
	})

	s.Engine.POST("/api/getbyliid", func(ctx *gin.Context) {
		s.router.GetByLIID(ctx)
	})

	s.Engine.POST("/api/getAllbyliid", func(ctx *gin.Context) {
		s.router.GetAllByLIID(ctx)
	})

	s.Engine.POST("/api/changewebhook", func(ctx *gin.Context) {
		s.router.ChangeWebhook(ctx)
	})

	s.Engine.POST("/api/getPersonalEmailByliid", func(ctx *gin.Context) {
		s.router.GetPersonalEmailByliid(ctx)
	})

	s.Engine.POST("/api/getProfessionalEmailsByliid", func(ctx *gin.Context) {
		s.router.GetProfessionalEmailsByliid(ctx)
	})
	
	s.Engine.POST("/api/test", func(ctx *gin.Context) {
		s.router.Test(ctx)
	})

	return s.Engine.Run(s.port)
}
