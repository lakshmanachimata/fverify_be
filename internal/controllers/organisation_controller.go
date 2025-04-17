package controllers

import (
	"net/http"

	"fverify_be/internal/models"
	"fverify_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Param organisation body models.OrganisationReq true "Organisation data"
// @Success 201 {object} models.Organisation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/organisations [post]
func (oc *OrganisationController) CreateOrganisation(c *gin.Context) {
	var reqOrg models.OrganisationReq
	if err := c.ShouldBindJSON(&reqOrg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var org models.Organisation
	org.OrgId = reqOrg.OrgId
	org.OrgName = reqOrg.OrgName
	org.Status = reqOrg.Status
	org.OrgUUID = uuid.New().String()
	// Generate a new UUID for the organisation
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
// @Param org_id path string true "Organisation ID"
// @Param organisation body models.OrganisationReq true "Updated organisation data"
// @Success 200 {object} models.Organisation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/organisations/{org_id} [put]
func (oc *OrganisationController) UpdateOrganisation(c *gin.Context) {
	org_id := c.Param("org_id")

	var org models.OrganisationReq
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the existing organisation to validate org_uuid
	existingOrg, err := oc.Service.GetOrganisationByID(c.Request.Context(), org_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
		return
	}

	existingOrg.OrgName = org.OrgName
	existingOrg.Status = org.Status
	existingOrg.OrgId = org.OrgId

	// Update the organisation
	err = oc.Service.UpdateOrganisation(c.Request.Context(), org_id, existingOrg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organisation"})
		return
	}

	// If the organisation status is updated to Inactive, update all users' status to Inactive
	if org.Status == models.OrgInActive {
		err = oc.Service.UpdateUsersStatusByOrgUUID(c.Request.Context(), existingOrg.OrgUUID, models.InActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update users' status"})
			return
		}
	}

	c.JSON(http.StatusOK, existingOrg)
}

// // DeleteOrganisation godoc
// // @Summary Delete an organisation
// // @Description Delete an organisation by its ID
// // @Tags Organisations
// // @Param X-API-Key header string true "API key"
// // @Param org_id path string true "Organisation ID"
// // @Success 204 "No Content"
// // @Failure 400 {object} ErrorResponse
// // @Failure 401 {object} InvalidAPIKeyResponse
// // @Failure 404 {object} NotFoundResponse
// // @Failure 500 {object} InternalErrorResponse
// // @Router /api/v1/organisations/{org_id} [delete]
func (oc *OrganisationController) DeleteOrganisation(c *gin.Context) {
	org_id := c.Param("org_id")

	err := oc.Service.DeleteOrganisation(c.Request.Context(), org_id)
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
// @Param X-API-Key header string true "API key"
// @Success 200 {array} models.Organisation
// @Failure 401 {object} InvalidAPIKeyResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/organisations [get]
func (oc *OrganisationController) GetAllOrganisations(c *gin.Context) {
	organisations, err := oc.Service.GetAllOrganisations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organisations"})
		return
	}

	c.JSON(http.StatusOK, organisations)
}
