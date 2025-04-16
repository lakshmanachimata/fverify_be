package controllers

import (
	"net/http"
	"strings"
	"time"

	"fverify_be/internal/models"
	"fverify_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
// @Param prospect body models.ProspecReqtModel true "Prospect data"
// @Success 201 {object} models.ProspectModel
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /prospects [post]
func (pc *ProspectController) CreateProspect(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(models.UserModel)
	var reqProspect models.ProspecReqtModel
	if err := c.ShouldBindJSON(&reqProspect); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map all fields from ProspecReqtModel to ProspectModel
	var prospect models.ProspectModel = models.ProspectModel{
		ProspectId:            reqProspect.ProspectId,
		ApplicantName:         reqProspect.ApplicantName,
		MobileNumber:          reqProspect.MobileNumber,
		Gender:                reqProspect.Gender,
		Age:                   reqProspect.Age,
		ResidentialAddress:    reqProspect.ResidentialAddress,
		YearsOfStay:           reqProspect.YearsOfStay,
		NumberOfFamilyMembers: reqProspect.NumberOfFamilyMembers,
		ReferenceName:         reqProspect.ReferenceName,
		ReferenceRelation:     reqProspect.ReferenceRelation,
		ReferenceMobile:       reqProspect.ReferenceMobile,
		EmploymentType:        reqProspect.EmploymentType,
		OfficeAddress:         reqProspect.OfficeAddress,
		YearsInCurrentOffice:  reqProspect.YearsInCurrentOffice,
		Role:                  reqProspect.Role,
		EmpId:                 reqProspect.EmpId,
		Status:                reqProspect.Status,
		PreviousExperience:    reqProspect.PreviousExperience,
		GrossSalary:           reqProspect.GrossSalary,
		NetSalary:             reqProspect.NetSalary,
		ColleagueName:         reqProspect.ColleagueName,
		ColleagueDesignation:  reqProspect.ColleagueDesignation,
		ColleagueMobile:       reqProspect.ColleagueMobile,
		UploadedImages:        reqProspect.UploadedImages,
		Remarks:               reqProspect.Remarks,
	}
	// Assign unique ID and timestamps
	prospect.UId = uuid.New().String()
	prospect.CreatedBy = authUser.Username
	prospect.CreatedTime = time.Now().UTC().Format(time.RFC3339) // Get current UTC time as string
	prospect.UpdatedBy = authUser.Username
	prospect.UpdatedTime = time.Now().UTC().Format(time.RFC3339) // Get current UTC time as string
	prospect.UpdateHistory = append(prospect.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join([]string{"Prospect created"}, ", "),
		UpdateBy:        authUser.Username,
	})
	// Call the service to create the prospect
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
