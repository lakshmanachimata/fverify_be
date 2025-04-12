package controllers

import (
	"net/http"
	"strconv"

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
// @Success 201 {object} SimpleResponse
// @Failure 400 {object} SimpleResponse
// @Failure 500 {object} SimpleResponse
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.Service.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} SimpleResponse
// @Failure 404 {object} SimpleResponse
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.Service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
