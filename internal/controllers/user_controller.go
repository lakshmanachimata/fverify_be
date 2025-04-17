package controllers

import (
	"net/http"
	"strings"
	"time"

	"fverify_be/internal/auth"
	"fverify_be/internal/models"
	"fverify_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	Service    *services.UserService
	OrgService *services.OrganisationService // Add this field
}

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

type InvalidAuthResponse struct {
	Error   string `json:"error" example:"Login Failed"`                    // Error message
	Details string `json:"details" example:"Invalid user name or password"` // Additional details about the error
}

type SuccessResponse struct {
	Message string `json:"message" example:"User created successfully"` // Success message
}

type InvalidAPIKeyResponse struct {
	Error   string `json:"error" example:"Invalid API key"`      // Error message
	Details string `json:"details" example:"API key is invalid"` // Additional details about the error
}

func NewUserController(userService *services.UserService, orgService *services.OrganisationService) *UserController {
	return &UserController{
		Service:    userService,
		OrgService: orgService, // Initialize OrgService
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param user body models.UserReqModel true "User data (all fields are mandatory)"
// @Success 201 {object} models.UserRespModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)
	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate orgUUID
	if reqUser.Org_Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgUUID is required"})
		return
	}

	// Check if the organisation exists and is active
	isActive, existingOrg := uc.OrgService.IsOrgActive(c.Request.Context(), reqUser.Org_Id)
	if existingOrg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation is inactive"})
		return
	}

	// Role-based access control
	switch authUser.Role {
	case string(models.Owner):
		// Owner can update all roles, no restrictions
	case string(models.Admin):
		if reqUser.Role != models.Admin &&
			reqUser.Role != models.OperationsLead &&
			reqUser.Role != models.FieldLead &&
			reqUser.Role != models.FieldExecutive &&
			reqUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admins can not create owner"})
			return
		}
	case string(models.OperationsLead):
		if reqUser.Role != models.OperationsLead &&
			reqUser.Role != models.FieldLead &&
			reqUser.Role != models.FieldExecutive &&
			reqUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operations Leads can not create owner / admin"})
			return
		}
	case string(models.OperationsExecutive):
		if reqUser.Role != models.FieldLead &&
			reqUser.Role != models.FieldExecutive &&
			reqUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operations Executives can not create owner / admin / operations lead"})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	var user models.UserModel
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgStatus = existingOrg.Status
	user.OrgUUID = authUser.OrgUUID
	user.CreatedTime = time.Now().UTC().Format(time.RFC3339)
	user.UpdatedTime = time.Now().UTC().Format(time.RFC3339)

	user.UpdateHistory = append(user.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join([]string{"User created"}, ", "),
		UpdateBy:        authUser.Username,
	})

	user.UId = uuid.New().String()

	createdUser, err := uc.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetByUserID godoc
// @Summary Get a user by userId
// @Description Retrieve a user by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param userId path int true "User ID"
// @Success 200 {object} models.UserRespModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/{userId} [get]
func (uc *UserController) GetUserByUserID(c *gin.Context) {
	idParam := c.Param("userId")
	user, err := uc.Service.GetByUserID(c.Request.Context(), idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Success 200 {array} models.UserRespModel
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.Service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUserByUId godoc
// @Summary Delete a user by uId
// @Description Delete a user by their unique uId
// @Tags Users
// @Param uId path int true "User uId"
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/uid/{uId} [delete]
// func (uc *UserController) DeleteUserByUId(c *gin.Context) {
// 	uIdParam := c.Param("uId")

// 	err = uc.Service.DeleteByUId(c.Request.Context(), uIdParam)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, ErrorResponse{
// 			Error:   "User not found",
// 			Details: "No user found with the given uId",
// 		})
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }

// DeleteUserByUserId godoc
// @Summary Delete a user by userId
// @Description Delete a user by their unique userId
// @Tags Users
// @Param userId path string true "User userId"
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/userid/{userId} [delete]
func (uc *UserController) DeleteUserByUserId(c *gin.Context) {
	userId := c.Param("userId")

	err := uc.Service.DeleteByUserId(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Details: "No user found with the given userId",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's details
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param uId path string true "User uId"
// @Param user body models.UserReqModel true "User data (all fields are mandatory)"
// @Success 200 {object} models.UserRespModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/uid/{uId} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)

	uIdParam := c.Param("uId")

	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the target user to validate roles
	targetUser, err := uc.Service.GetByUserUID(c.Request.Context(), uIdParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Role-based access control
	switch authUser.Role {
	case string(models.Owner):
		// Owner can update all roles, no restrictions
	case string(models.Admin):
		if targetUser.Role != models.Admin &&
			targetUser.Role != models.OperationsLead &&
			targetUser.Role != models.FieldLead &&
			targetUser.Role != models.FieldExecutive &&
			targetUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admins can only update specific roles"})
			return
		}
	case string(models.OperationsLead):
		if targetUser.Role != models.OperationsLead &&
			targetUser.Role != models.FieldLead &&
			targetUser.Role != models.FieldExecutive &&
			targetUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operations Leads can only update specific roles"})
			return
		}
	case string(models.OperationsExecutive):
		if targetUser.Role != models.FieldLead &&
			targetUser.Role != models.FieldExecutive &&
			targetUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operations Executives can only update specific roles"})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update user data
	var user models.UserModel
	user.UId = uIdParam
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgUUID = authUser.OrgUUID

	uUser, err := uc.Service.UpdateUser(c.Request.Context(), &user, authUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, uUser)
}

// LoginUser godoc
// @Summary Login a user
// @Description Validate username and password, and return user details with a token
// @Tags Users
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/login [post]
func (uc *UserController) LoginUser(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the organisation exists and is active
	isActive, existingOrg := uc.OrgService.IsOrgActive(c.Request.Context(), loginRequest.OrgId)
	if existingOrg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid organisation"})
		return
	}

	// Validate the user
	user, err := uc.Service.LoginUser(c.Request.Context(), loginRequest.Username, loginRequest.Password, existingOrg.OrgUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if user.Status == models.InActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is inactive"})
		return
	}
	// Update user status to Active
	if user.Status != models.Active {
		err = uc.Service.UpdateUserStatus(c.Request.Context(), user.UserId, string(models.Active))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
			return
		}
		user.Status = models.Active // Update the status in the user object
	}

	// Generate the token
	token, err := auth.GenerateAuthToken(user.UserId, user.Username, user.UId, string(user.Role), string(user.Status), user.MobileNumber, user.OrgUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with user details and token
	c.JSON(http.StatusOK, models.LoginResponse{
		UId:          user.UId,
		UserId:       user.UserId,
		Username:     user.Username,
		Role:         string(user.Role),
		Status:       string(user.Status),
		MobileNumber: user.MobileNumber,
		Token:        token,
	})
}

// SetPassword godoc
// @Summary Set a user's password
// @Description Set a new password for a user
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param org_id  header string true "Organisation Id"
// @Param uId path int true "User uId"
// @Param password body models.SetPasswordRequest true "New password"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/uid/{uId}/setpassword [put]
func (uc *UserController) SetPassword(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)

	uIdParam := c.Param("uId")

	var request models.SetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the target user to validate roles
	targetUser, err := uc.Service.GetByUserUID(c.Request.Context(), uIdParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Role-based access control
	if authUser.Role == string(models.Admin) || authUser.Role == string(models.Owner) {
		if targetUser.Role != models.Admin && targetUser.Role != models.OperationsLead {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admins and Owners can only set passwords for Admins and Operations Leads"})
			return
		}
	} else if authUser.Role == string(models.OperationsLead) {
		if targetUser.Role != models.OperationsLead && targetUser.Role != models.FieldLead &&
			targetUser.Role != models.FieldExecutive && targetUser.Role != models.OperationsExecutive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Operations Leads can only set passwords for specific roles"})
			return
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update the password
	err = uc.Service.SetPassword(c.Request.Context(), uIdParam, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// CreateAdminUser godoc
// @Summary Create a new admin user
// @Description Create a new admin user in the system (requires API key)
// @Tags Users
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param user body models.UserReqModel true "Admin user data"
// @Success 201 {object} models.UserRespModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/admin/create [post]
func (uc *UserController) CreateAdmin(c *gin.Context) {
	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate orgUUID
	if reqUser.Org_Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgUUID is required"})
		return
	}

	// Check if the organisation exists and is active
	isActive, existingOrg := uc.OrgService.IsOrgActive(c.Request.Context(), reqUser.Org_Id)
	if existingOrg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation is inactive"})
		return
	}

	var user models.UserModel
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgStatus = existingOrg.Status
	user.OrgUUID = existingOrg.OrgUUID
	user.UId = uuid.New().String()

	user.CreatedTime = time.Now().UTC().Format(time.RFC3339)
	user.UpdatedTime = time.Now().UTC().Format(time.RFC3339)

	user.UpdateHistory = append(user.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join([]string{"Admin created"}, ", "),
		UpdateBy:        "System",
	})
	// Ensure the role is set to Admin
	user.Role = models.Admin

	createdUser, err := uc.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// CreateAdminUser godoc
// @Summary Create a new owner
// @Description Create a new admin user in the system (requires API key)
// @Tags Users
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param user body models.UserReqModel true "User data (all fields are mandatory)"
// @Success 201 {object} models.UserRespModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/owner/create [post]
func (uc *UserController) CreateOwner(c *gin.Context) {
	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate orgUUID
	if reqUser.Org_Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgUUID is required"})
		return
	}

	// Check if the organisation exists and is active
	isActive, existingOrg := uc.OrgService.IsOrgActive(c.Request.Context(), reqUser.Org_Id)
	if existingOrg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation is inactive"})
		return
	}

	var user models.UserModel
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgStatus = existingOrg.Status
	user.OrgUUID = existingOrg.OrgUUID
	user.UId = uuid.New().String()
	user.CreatedTime = time.Now().UTC().Format(time.RFC3339)
	user.UpdatedTime = time.Now().UTC().Format(time.RFC3339)
	// Set the update history

	user.UpdateHistory = append(user.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join([]string{"Owner created"}, ", "),
		UpdateBy:        "System",
	})
	// Ensure the role is set to Admin
	user.Role = models.Owner

	createdUser, err := uc.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Retrieve all user roles for a given organisation
// @Tags Users
// @Accept json
// @Produce json
// @Param org_id  header string true "Organisation Id"
// @Success 200 {array} string
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/users/roles [get]
func (uc *UserController) GetUserRoles(c *gin.Context) {
	org_id := c.GetHeader("org_id")
	if org_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id is required"})
		return
	}

	// Check if the organisation exists and is active
	isActive, existingOrg := uc.OrgService.IsOrgActive(c.Request.Context(), org_id)
	if existingOrg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found or inactive"})
		return
	}

	// Return the list of user roles
	roles := []string{
		string(models.Admin),
		string(models.OperationsLead),
		string(models.FieldLead),
		string(models.FieldExecutive),
		string(models.Owner),
		string(models.OperationsExecutive),
	}
	c.JSON(http.StatusOK, roles)
}
