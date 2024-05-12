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

	s.Engine.POST("/scandb", func(ctx *gin.Context) {
		s.router.Pep1(ctx)
	})

	s.Engine.POST("/getPersonalEmail", func(ctx *gin.Context) {
		s.router.GetPersonalEmail(ctx)
	})

	s.Engine.POST("/getProfessionalEmails", func(ctx *gin.Context) {
		s.router.GetProfessionalEmails(ctx)
	})
	s.Engine.POST("/getBothEmails", func(ctx *gin.Context) {
		s.router.GetBothEmails(ctx)
	})

	s.Engine.POST("/getbyliid", func(ctx *gin.Context) {
		s.router.GetByLIID(ctx)
	})

	s.Engine.POST("/getAllbyliid", func(ctx *gin.Context) {
		s.router.GetAllByLIID(ctx)
	})

	s.Engine.POST("/changewebhook", func(ctx *gin.Context) {
		s.router.ChangeWebhook(ctx)
	})

	s.Engine.POST("/getPersonalEmailByliid", func(ctx *gin.Context) {
		s.router.GetPersonalEmailByliid(ctx)
	})

	s.Engine.POST("/getProfessionalEmailsByliid", func(ctx *gin.Context) {
		s.router.GetProfessionalEmailsByliid(ctx)
	})

	return s.Engine.Run(s.port)
}
