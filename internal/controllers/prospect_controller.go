package controllers

import (
	"net/http"

	"kowtha_be/internal/models"
	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`          // Error message
	Details string `json:"details" example:"Invalid input data"` // Additional details about the error
}

type InternalErrorResponse struct {
	Error   string `json:"error" example:"Internal Server Error"` // Error message
	Details string `json:"details" example:"Server Error"`        // Additional details about the error
}

type NotFoundResponse struct {
	Error   string `json:"error" example:"No data"`                                // Error message
	Details string `json:"details" example:"No Data found for the input provided"` // Additional details about the error
}

type ProspectController struct {
	Service *services.ProspectService
}

func NewProspectController(service *services.ProspectService) *ProspectController {
	return &ProspectController{Service: service}
}

// CreateProspect godoc
// @Summary Create a new prospect
// @Description Create a new prospect in the system
// @Tags Prospects
// @Accept json
// @Produce json
// @Param prospect body models.ProspectModel true "Prospect data"
// @Success 201 {object} models.ProspectModel
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /prospects [post]
func (pc *ProspectController) CreateProspect(c *gin.Context) {
	var prospect models.ProspectModel
	if err := c.ShouldBindJSON(&prospect); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.CreateProspect(c.Request.Context(), &prospect); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prospect"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Prospect created"})
}

// GetProspect godoc
// @Summary Get a prospect by ID
// @Description Retrieve a prospect by their unique ID
// @Tags Prospects
// @Accept json
// @Produce json
// @Param id path string true "Prospect ID"
// @Success 200 {object} models.ProspectModel
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /prospects/{id} [get]
func (pc *ProspectController) GetProspect(c *gin.Context) {
	id := c.Param("id")
	prospect, err := pc.Service.GetProspectByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prospect not found"})
		return
	}

	c.JSON(http.StatusOK, prospect)
}
