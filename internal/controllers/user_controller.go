package controllers

import (
	"net/http"

	"kowtha_be/internal/models"
	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserModel true "User data"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} InternalErrorResponse
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure uId is not set by the payload
	user.UId = 0

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
// @Param userId path int true "User ID"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
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
// @Success 200 {array} models.UserModel
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
