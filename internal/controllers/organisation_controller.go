package controllers

import (
	"net/http"

	"kowtha_be/internal/models"
	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
)

type OrganisationController struct {
	Service *services.OrganisationService
}

func NewOrganisationController(service *services.OrganisationService) *OrganisationController {
	return &OrganisationController{Service: service}
}

// CreateOrganisation godoc
// @Summary Create a new organisation
// @Description Create a new organisation in the system
// @Tags Organisations
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param organisation body models.Organisation true "Organisation data"
// @Success 201 {object} models.Organisation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /organisations [post]
func (oc *OrganisationController) CreateOrganisation(c *gin.Context) {
	var org models.Organisation
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdOrg, err := oc.Service.CreateOrganisation(c.Request.Context(), &org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organisation"})
		return
	}

	c.JSON(http.StatusCreated, createdOrg)
}

// UpdateOrganisation godoc
// @Summary Update an organisation
// @Description Update an existing organisation's details
// @Tags Organisations
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param orgId path string true "Organisation ID"
// @Param organisation body models.Organisation true "Updated organisation data"
// @Success 200 {object} models.Organisation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /organisations/{orgId} [put]
func (oc *OrganisationController) UpdateOrganisation(c *gin.Context) {
	orgId := c.Param("orgId")

	var org models.Organisation
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := oc.Service.UpdateOrganisation(c.Request.Context(), orgId, &org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organisation"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganisation godoc
// @Summary Delete an organisation
// @Description Delete an organisation by its ID
// @Tags Organisations
// @Param X-API-Key header string true "API key"
// @Param orgId path string true "Organisation ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /organisations/{orgId} [delete]
func (oc *OrganisationController) DeleteOrganisation(c *gin.Context) {
	orgId := c.Param("orgId")

	err := oc.Service.DeleteOrganisation(c.Request.Context(), orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organisation"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllOrganisations godoc
// @Summary Get all organisations
// @Description Retrieve all organisations in the system
// @Tags Organisations
// @Accept json
// @Produce json
// @Success 200 {array} models.Organisation
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /organisations [get]
func (oc *OrganisationController) GetAllOrganisations(c *gin.Context) {
	organisations, err := oc.Service.GetAllOrganisations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organisations"})
		return
	}

	c.JSON(http.StatusOK, organisations)
}
