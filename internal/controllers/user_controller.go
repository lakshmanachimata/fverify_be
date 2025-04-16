package controllers

import (
	"net/http"
	"strconv"

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
// @Param user body models.UserReqModel true "User data"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.UserModel
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.CreatedTime = reqUser.CreatedTime
	user.UpdatedTime = reqUser.UpdatedTime
	user.UpdateHistory = reqUser.UpdateHistory
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgStatus = reqUser.OrgStatus
	user.OrgUUID = reqUser.OrgUUID

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
// @Param userId path int true "User ID"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users/{userId} [get]
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
// @Success 200 {array} models.UserModel
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
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
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/uid/{uId} [delete]
func (uc *UserController) DeleteUserByUId(c *gin.Context) {
	uIdParam := c.Param("uId")
	uId, err := strconv.Atoi(uIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid uId",
			Details: "uId must be a valid integer",
		})
		return
	}

	err = uc.Service.DeleteByUId(c.Request.Context(), uId)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Details: "No user found with the given uId",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteUserByUserId godoc
// @Summary Delete a user by userId
// @Description Delete a user by their unique userId
// @Tags Users
// @Param userId path string true "User userId"
// @Param Authorization header string true "Bearer token"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users/userid/{userId} [delete]
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
// @Param uId path string true "User uId"
// @Param user body models.UserModel true "Updated user data"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users/uid/{uId} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)

	uIdParam := c.Param("uId")

	var reqUser models.UserReqModel
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.UserModel
	user.UserId = reqUser.UserId
	user.Username = reqUser.Username
	user.Password = reqUser.Password
	user.Role = reqUser.Role
	user.Status = reqUser.Status
	user.CreatedTime = reqUser.CreatedTime
	user.UpdatedTime = reqUser.UpdatedTime
	user.UpdateHistory = reqUser.UpdateHistory
	user.Remarks = reqUser.Remarks
	user.MobileNumber = reqUser.MobileNumber
	user.OrgStatus = reqUser.OrgStatus
	user.OrgUUID = reqUser.OrgUUID
	user.UId = uIdParam

	// Restrict role changes for Operations Lead
	if authUser.Role == string(models.OperationsLead) && user.Role != "" && string(user.Role) != authUser.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": "Operations Lead cannot change roles"})
		return
	}

	err := uc.Service.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
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
// @Router /users/login [post]
func (uc *UserController) LoginUser(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the organisation exists and is active
	isActive, err := uc.OrgService.IsOrgActive(c.Request.Context(), loginRequest.OrgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate organisation"})
		return
	}
	if !isActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid organisation"})
		return
	}

	// Validate the user
	user, err := uc.Service.LoginUser(c.Request.Context(), loginRequest.Username, loginRequest.Password, loginRequest.OrgId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
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
// @Param uId path int true "User uId"
// @Param password body models.SetPasswordRequest true "New password"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users/uid/{uId}/setpassword [put]
func (uc *UserController) SetPassword(c *gin.Context) {
	claims, _ := c.Get("user")
	authUser := claims.(*auth.AuthTokenClaims)

	uIdParam := c.Param("uId")
	uId, err := strconv.Atoi(uIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uId"})
		return
	}

	var request models.SetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the target user to validate roles
	targetUser, err := uc.Service.GetByUserUID(c.Request.Context(), uId)
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
	err = uc.Service.SetPassword(c.Request.Context(), uId, request.Password)
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
// @Param user body models.UserModel true "Admin user data"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /admin/create [post]
func (uc *UserController) CreateAdmin(c *gin.Context) {
	var user models.UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
// @Param user body models.UserModel true "Admin user data"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /owner/create [post]
func (uc *UserController) CreateOwner(c *gin.Context) {
	var user models.UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
// @Param orgId query string true "Organisation ID"
// @Success 200 {array} string
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 401 {object} InvalidAuthResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users/roles [get]
func (uc *UserController) GetUserRoles(c *gin.Context) {
	orgId := c.Query("orgId")
	if orgId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId is required"})
		return
	}

	// Check if the organisation exists and is active
	isActive, err := uc.OrgService.IsOrgActive(c.Request.Context(), orgId)
	if err != nil {
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
