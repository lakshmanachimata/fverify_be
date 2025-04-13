package controllers

import (
	"net/http"
	"strconv"

	"kowtha_be/internal/models"
	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

// DeleteUserByUId godoc
// @Summary Delete a user by uId
// @Description Delete a user by their unique uId
// @Tags Users
// @Param uId path int true "User uId"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
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
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Param uId path int true "User uId"
// @Param user body models.UserModel true "Updated user data"
// @Success 200 {object} models.UserModel
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/uid/{uId} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	uIdParam := c.Param("uId")
	uId, err := strconv.Atoi(uIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid uId",
			Details: "uId must be a valid integer",
		})
		return
	}

	var user models.UserModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input",
			Details: err.Error(),
		})
		return
	}

	// Ensure the uId in the payload matches the path parameter
	user.UId = uId

	err = uc.Service.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Details: "No user found with the given uId",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
