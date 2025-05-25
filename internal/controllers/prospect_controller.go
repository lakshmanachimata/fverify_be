package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"fverify_be/internal/auth"
	"fverify_be/internal/models"
	"fverify_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProspectCountMessage struct {
	Count int `json:"count"`
}

type ProspectController struct {
	Service *services.ProspectService
}

func NewProspectController(service *services.ProspectService) *ProspectController {
	return &ProspectController{Service: service}
}

// GetProspectsCount godoc
// @Summary Get total count of prospects
// @Description Retrieve the total count of prospects in the system
// @Tags Prospects
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id header string true "Organisation Id"
// @Success 200 {object} ProspectCountMessage
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/prospects/count [get]
func (pc *ProspectController) GetProspectsCount(c *gin.Context) {
	// claims, _ := c.Get("user")
	// authUser := claims.(*auth.AuthTokenClaims)
	// Call the service to get the total count of prospects
	count, err := pc.Service.GetProspectsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve prospects count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetProspects godoc
// @Summary Get a list of prospects
// @Description Retrieve a list of prospects with pagination using skip and limit values
// @Tags Prospects
// @Accept json
// @Produce json
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Number of records to retrieve" default(10)
// @Param Authorization header string true "Bearer token"
// @Param org_id header string true "Organisation Id"
// @Success 200 {array} models.Prospect
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/prospects [get]
func (pc *ProspectController) GetProspects(c *gin.Context) {
	// claims, _ := c.Get("user")
	// authUser := claims.(*auth.AuthTokenClaims)
	// Parse query parameters
	skip := 0
	limit := 10

	if s := c.Query("skip"); s != "" {
		if parsedSkip, err := strconv.Atoi(s); err == nil {
			skip = parsedSkip
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skip value"})
			return
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
	}

	// Call the service to get prospects
	prospects, err := pc.Service.GetProspects(c.Request.Context(), skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve prospects"})
		return
	}

	c.JSON(http.StatusOK, prospects)
}

// CreateProspect godoc
// @Summary Create a new prospect
// @Description Create a new prospect in the system
// @Tags Prospects
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param prospect body models.ProspecReq true "Prospect data"
// @Success 201 {object} models.Prospect
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/prospects [post]
func (pc *ProspectController) CreateProspect(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)
	var reqProspect models.ProspecReq
	if err := c.ShouldBindJSON(&reqProspect); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map all fields from ProspecReq to Prospect
	var prospect models.Prospect = models.Prospect{
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
	prospect.NameVerified = false
	prospect.MobileVerified = false
	prospect.ResAddressVerified = false
	prospect.OffAddressVerified = false
	prospect.RoleVerified = false
	prospect.EmpIdVerified = false
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

	c.JSON(http.StatusCreated, prospect)
}

// GetProspect godoc
// @Summary Get a prospect by ID
// @Description Retrieve a prospect by their unique ID
// @Tags Prospects
// @Accept json
// @Produce json
// @Param uid path string true "Prospect UID"
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Success 200 {object} models.Prospect
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/prospects/{id} [get]
func (pc *ProspectController) GetProspect(c *gin.Context) {
	uid := c.Param("uid")
	// claims, _ := c.Get("user")
	// authUser := claims.(*auth.AuthTokenClaims)
	prospect, err := pc.Service.GetProspectByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prospect not found"})
		return
	}

	c.JSON(http.StatusOK, prospect)
}

// UpdateProspect godoc
// @Summary Update an existing prospect
// @Description Update an existing prospect in the system. Update comments are generated based on differences from the earlier prospect state.
// @Tags Prospects
// @Accept json
// @Produce json
// @Param uid path string true "Prospect UId"
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param prospect body models.ProspecReq true "Updated prospect data"
// @Success 200 {object} models.Prospect
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/prospects/{uid} [put]
func (pc *ProspectController) UpdateProspect(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)
	uId := c.Param("uid")

	// Fetch the existing prospect
	existingProspect, err := pc.Service.GetProspectByID(c.Request.Context(), uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prospect not found"})
		return
	}

	var reqProspect models.ProspecReq
	if err := c.ShouldBindJSON(&reqProspect); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate update comments by comparing the existing and new prospect data
	// Generate update comments by comparing the existing and new prospect data
	var updateComments []string
	if existingProspect.ProspectId != reqProspect.ProspectId && reqProspect.ProspectId != "" {
		updateComments = append(updateComments, "ProspectId updated")
	}
	if existingProspect.ApplicantName != reqProspect.ApplicantName && reqProspect.ProspectId != "" {
		updateComments = append(updateComments, "ApplicantName updated")
	}
	if existingProspect.MobileNumber != reqProspect.MobileNumber && reqProspect.ProspectId != "" {
		updateComments = append(updateComments, "MobileNumber updated")
	}
	if existingProspect.Gender != reqProspect.Gender && reqProspect.Gender != "" {
		updateComments = append(updateComments, "Gender updated")
	}
	if existingProspect.Age != reqProspect.Age && reqProspect.Age > 0 {
		updateComments = append(updateComments, "Age updated")
	}
	if existingProspect.ResidentialAddress != reqProspect.ResidentialAddress && reqProspect.ResidentialAddress != "" {
		updateComments = append(updateComments, "ResidentialAddress updated")
	}
	if existingProspect.YearsOfStay != reqProspect.YearsOfStay && reqProspect.YearsOfStay > 0 {
		updateComments = append(updateComments, "YearsOfStay updated")
	}
	if existingProspect.NumberOfFamilyMembers != reqProspect.NumberOfFamilyMembers && reqProspect.NumberOfFamilyMembers > 0 {
		updateComments = append(updateComments, "NumberOfFamilyMembers updated")
	}
	if existingProspect.ReferenceName != reqProspect.ReferenceName && reqProspect.ReferenceName != "" {
		updateComments = append(updateComments, "ReferenceName updated")
	}
	if existingProspect.ReferenceRelation != reqProspect.ReferenceRelation && reqProspect.ReferenceRelation != "" {
		updateComments = append(updateComments, "ReferenceRelation updated")
	}
	if existingProspect.ReferenceMobile != reqProspect.ReferenceMobile && reqProspect.ReferenceMobile != "" {
		updateComments = append(updateComments, "ReferenceMobile updated")
	}
	if existingProspect.EmploymentType != reqProspect.EmploymentType && reqProspect.EmploymentType != "" {
		updateComments = append(updateComments, "EmploymentType updated")
	}
	if existingProspect.OfficeAddress != reqProspect.OfficeAddress && reqProspect.OfficeAddress != "" {
		updateComments = append(updateComments, "OfficeAddress updated")
	}
	if existingProspect.YearsInCurrentOffice != reqProspect.YearsInCurrentOffice && reqProspect.YearsInCurrentOffice > 0 {
		updateComments = append(updateComments, "YearsInCurrentOffice updated")
	}
	if existingProspect.Role != reqProspect.Role && reqProspect.Role != "" {
		updateComments = append(updateComments, "Role updated")
	}
	if existingProspect.EmpId != reqProspect.EmpId && reqProspect.EmpId != "" {
		updateComments = append(updateComments, "EmpId updated")
	}
	if existingProspect.Status != reqProspect.Status && reqProspect.Status != "" {
		updateComments = append(updateComments, "Status updated")
	}
	if existingProspect.PreviousExperience != reqProspect.PreviousExperience {
		updateComments = append(updateComments, "PreviousExperience updated")
	}
	if existingProspect.GrossSalary != reqProspect.GrossSalary && reqProspect.GrossSalary > 0 {
		updateComments = append(updateComments, "GrossSalary updated")
	}
	if existingProspect.NetSalary != reqProspect.NetSalary && reqProspect.NetSalary > 0 {
		updateComments = append(updateComments, "NetSalary updated")
	}
	if existingProspect.ColleagueName != reqProspect.ColleagueName && reqProspect.ColleagueName != "" {
		updateComments = append(updateComments, "ColleagueName updated")
	}
	if existingProspect.ColleagueDesignation != reqProspect.ColleagueDesignation && reqProspect.ColleagueDesignation != "" {
		updateComments = append(updateComments, "ColleagueDesignation updated")
	}
	if existingProspect.ColleagueMobile != reqProspect.ColleagueMobile && reqProspect.ColleagueMobile != "" {
		updateComments = append(updateComments, "ColleagueMobile updated")
	}
	if existingProspect.Remarks != reqProspect.Remarks && reqProspect.Remarks != "" {
		updateComments = append(updateComments, "Remarks updated")
	}
	if existingProspect.NameVerified != reqProspect.NameVerified {
		updateComments = append(updateComments, "NameVerified updated")
	}
	if existingProspect.MobileVerified != reqProspect.MobileVerified {
		updateComments = append(updateComments, "MobileVerified updated")
	}
	if existingProspect.ResAddressVerified != reqProspect.ResAddressVerified {
		updateComments = append(updateComments, "Residential Address updated")
	}
	if existingProspect.OffAddressVerified != reqProspect.OffAddressVerified {
		updateComments = append(updateComments, "Office Address updated")
	}
	if existingProspect.RoleVerified != reqProspect.RoleVerified {
		updateComments = append(updateComments, "Role / Business updated")
	}
	if existingProspect.ResAddressVerified != reqProspect.ResAddressVerified {
		updateComments = append(updateComments, "Employee Id / Business Id updated")
	}

	// Map updated fields from ProspecReq to Prospect
	existingProspect.ProspectId = reqProspect.ProspectId
	existingProspect.ApplicantName = reqProspect.ApplicantName
	existingProspect.MobileNumber = reqProspect.MobileNumber
	existingProspect.Gender = reqProspect.Gender
	existingProspect.Age = reqProspect.Age
	existingProspect.ResidentialAddress = reqProspect.ResidentialAddress
	existingProspect.YearsOfStay = reqProspect.YearsOfStay
	existingProspect.NumberOfFamilyMembers = reqProspect.NumberOfFamilyMembers
	existingProspect.ReferenceName = reqProspect.ReferenceName
	existingProspect.ReferenceRelation = reqProspect.ReferenceRelation
	existingProspect.ReferenceMobile = reqProspect.ReferenceMobile
	existingProspect.EmploymentType = reqProspect.EmploymentType
	existingProspect.OfficeAddress = reqProspect.OfficeAddress
	existingProspect.YearsInCurrentOffice = reqProspect.YearsInCurrentOffice
	existingProspect.Role = reqProspect.Role
	existingProspect.EmpId = reqProspect.EmpId
	existingProspect.Status = reqProspect.Status
	existingProspect.PreviousExperience = reqProspect.PreviousExperience
	existingProspect.GrossSalary = reqProspect.GrossSalary
	existingProspect.NetSalary = reqProspect.NetSalary
	existingProspect.ColleagueName = reqProspect.ColleagueName
	existingProspect.ColleagueDesignation = reqProspect.ColleagueDesignation
	existingProspect.ColleagueMobile = reqProspect.ColleagueMobile
	existingProspect.UploadedImages = reqProspect.UploadedImages
	existingProspect.Remarks = reqProspect.Remarks
	existingProspect.NameVerified = reqProspect.NameVerified
	existingProspect.MobileVerified = reqProspect.MobileVerified
	existingProspect.ResAddressVerified = reqProspect.ResAddressVerified
	existingProspect.OffAddressVerified = reqProspect.OffAddressVerified
	existingProspect.RoleVerified = reqProspect.RoleVerified
	existingProspect.EmpIdVerified = reqProspect.EmpIdVerified

	// Update timestamps and history
	existingProspect.UpdatedBy = authUser.Username
	existingProspect.UpdatedTime = time.Now().UTC().Format(time.RFC3339)
	existingProspect.UpdateHistory = append(existingProspect.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join(updateComments, ", "),
		UpdateBy:        authUser.Username,
	})

	// Call the service to update the prospect
	if err := pc.Service.UpdateProspect(c.Request.Context(), existingProspect); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prospect"})
		return
	}

	c.JSON(http.StatusOK, existingProspect)
}
