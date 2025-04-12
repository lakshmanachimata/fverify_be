package controllers

import (
	"net/http"

	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
)

type ProspectController struct {
	Service *services.ProspectService
}

func NewProspectController(service *services.ProspectService) *ProspectController {
	return &ProspectController{Service: service}
}

func (pc *ProspectController) CreateProspect(c *gin.Context) {
	// Handle creating a prospect
	c.JSON(http.StatusCreated, gin.H{"message": "Prospect created"})
}

func (pc *ProspectController) GetProspect(c *gin.Context) {
	// Handle retrieving a prospect
	c.JSON(http.StatusOK, gin.H{"message": "Prospect retrieved"})
}
