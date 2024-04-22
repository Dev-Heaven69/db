package router

import (
	"log"
	"net/http"

	"github.com/DevHeaven/db/domain/models"
	"github.com/DevHeaven/db/internal/logic"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Logic logic.Logic
}

func ProvideRouter(l logic.Logic) Router {
	return Router{l}
}

func (r Router) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server healthy"})
}

func (r Router) Pep1(c *gin.Context) {
	var req models.Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := r.Logic.FindPep1(file, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("response: ", resp)

	c.JSON(http.StatusOK, resp)
}
