package main

import (
	"context"
	"log"

	"kowtha_be/internal/auth"
	"kowtha_be/internal/controllers"
	"kowtha_be/internal/repositories"
	"kowtha_be/internal/services"

	"kowtha_be/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// @title Kowtha API
// @version 1.0
// @description This is the API documentation for the Kowtha backend.
// @host localhost:9000
// @BasePath /
// @schemes http
func main() {
	// Set up MongoDB connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://kowtha:nearhop%40123@applicants.cq0no3x.mongodb.net/?appName=Applicants").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	// Initialize repositories
	prospectRepo := repositories.NewProspectRepository(client, "kowtha_db", "prospects")
	userRepo := repositories.NewUserRepository(client, "kowtha_db", "users")
	orgRepo := repositories.NewOrganisationRepository(client, "kowtha_db", "orgs")

	// Initialize services
	prospectService := services.NewProspectService(prospectRepo)
	userService := services.NewUserService(userRepo)
	orgService := services.NewOrganisationService(orgRepo)

	// Initialize controllers
	prospectController := controllers.NewProspectController(prospectService)
	userController := controllers.NewUserController(userService)
	organisationController := controllers.NewOrganisationController(orgService)

	// Set up Gin router
	router := gin.Default()

	// Swagger setup
	docs.SwaggerInfo.Title = "Kowtha API"
	docs.SwaggerInfo.Description = "This is the API documentation for the Kowtha backend."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:9000"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// @Summary Create a new organisation
	// @Description Create a new organisation in the system
	// @Tags Organisations
	// @Accept json
	// @Produce json
	// @Param X-API-Key header string true "API key"
	// @Param organisation body models.Organisation true "Organisation data"
	// @Success 201 {object} models.Organisation
	// @Failure 400 {object} gin.H{"error": "Invalid input"}
	// @Failure 401 {object} gin.H{"error": "Invalid API key"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /organisations [post]
	router.POST("/organisations", auth.OrgAPIKeyMiddleware(), organisationController.CreateOrganisation)
	// @Summary Update an organisation
	// @Description Update an existing organisation's details
	// @Tags Organisations
	// @Accept json
	// @Produce json
	// @Param X-API-Key header string true "API key"
	// @Param orgId path string true "Organisation ID"
	// @Param organisation body models.Organisation true "Updated organisation data"
	// @Success 200 {object} models.Organisation
	// @Failure 400 {object} gin.H{"error": "Invalid input"}
	// @Failure 401 {object} gin.H{"error": "Invalid API key"}
	// @Failure 404 {object} gin.H{"error": "Organisation not found"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /organisations/{orgId} [put]
	router.PUT("/organisations/:orgId", auth.OrgAPIKeyMiddleware(), organisationController.UpdateOrganisation)
	// @Summary Delete an organisation
	// @Description Delete an organisation by its ID
	// @Tags Organisations
	// @Param X-API-Key header string true "API key"
	// @Param orgId path string true "Organisation ID"
	// @Success 204 "No Content"
	// @Failure 401 {object} gin.H{"error": "Invalid API key"}
	// @Failure 404 {object} gin.H{"error": "Organisation not found"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /organisations/{orgId} [delete]
	router.DELETE("/organisations/:orgId", auth.OrgAPIKeyMiddleware(), organisationController.DeleteOrganisation)

	// @Summary Get all organisations
	// @Description Retrieve all organisations in the system
	// @Tags Organisations
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.Organisation
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /organisations [get]
	router.GET("/organisations", auth.GetOrgAPIKeyMiddleware(), organisationController.GetAllOrganisations)

	// @Summary Login a user
	// @Description Validate username and password, and return user details with a token
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param login body models.LoginRequest true "Login credentials"
	// @Success 200 {object} models.LoginResponse
	// @Failure 401 {object} gin.H{"error": "Invalid username or password"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users/login [post]
	router.POST("/users/login", userController.LoginUser)

	// User APIs
	// @Summary Create a new user
	// @Description Create a new user in the system
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Accept json
	// @Produce json
	// @Param user body models.UserModel true "User data"
	// @Success 201 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Bad Request"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users [post]
	router.POST("/users", auth.AuthMiddleware("Admin", "Owner"), userController.CreateUser)

	// @Summary Update a user
	// @Description Update an existing user's details
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Accept json
	// @Produce json
	// @Param uId path int true "User uId"
	// @Param user body models.UserModel true "Updated user data"
	// @Success 200 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Invalid uId"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users/uid/{uId} [put]
	router.PUT("/users/uid/:uId", auth.AuthMiddleware("Admin", "Owner", "Operations Lead"), userController.UpdateUser)

	// @Summary Get all users
	// @Description Retrieve all users in the system
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.UserModel
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users [get]
	router.GET("/users", auth.AuthMiddleware("Admin", "Owner", "Operations Lead", "Operations Executive"), userController.GetAllUsers)

	// @Summary Get a user by ID
	// @Description Retrieve a user by their unique ID
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Accept json
	// @Produce json
	// @Param userId path string true "User ID"
	// @Success 200 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Router /users/{userId} [get]
	router.GET("/users/:userId", auth.AuthMiddleware("Admin", "Owner", "Operations Lead", "Operations Executive"), userController.GetUserByUserID)

	// @Summary Delete a user by uId
	// @Description Delete a user by their unique uId
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Param uId path int true "User uId"
	// @Success 204 "No Content"
	// @Failure 400 {object} gin.H{"error": "Invalid uId"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Router /users/uid/{uId} [delete]
	router.DELETE("/users/uid/:uId", auth.AuthMiddleware("Admin", "Owner"), userController.DeleteUserByUId)

	// @Summary Delete a user by userId
	// @Description Delete a user by their unique userId
	// @Tags Users
	// @Param Authorization header string true "Bearer token"
	// @Param userId path string true "User userId"
	// @Success 204 "No Content"
	// @Failure 400 {object} gin.H{"error": "Invalid userId"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Router /users/userid/{userId} [delete]
	router.DELETE("/users/userid/:userId", auth.AuthMiddleware("Admin", "Owner"), userController.DeleteUserByUserId)

	// @Summary Set a user's password
	// @Description Set a new password for a user
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param Authorization header string true "Bearer token"
	// @Param uId path int true "User uId"
	// @Param password body models.SetPasswordRequest true "New password"
	// @Success 200 {object} gin.H{"message": "Password updated successfully"}
	// @Failure 400 {object} gin.H{"error": "Invalid input"}
	// @Failure 403 {object} gin.H{"error": "Access denied"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users/uid/{uId}/setpassword [put]
	router.PUT("/users/uid/:uId/setpassword", auth.AuthMiddleware("Admin", "Owner", "Operations Lead"), userController.SetPassword)

	// @Summary Create a new admin user
	// @Description Create a new admin user in the system (requires API key)
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param X-API-Key header string true "API key"
	// @Param user body models.UserModel true "Admin user data"
	// @Success 201 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Invalid input"}
	// @Failure 401 {object} gin.H{"error": "Invalid API key"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /admin/create [post]
	router.POST("/users/admin/create", auth.APIKeyMiddleware(), userController.CreateAdmin)

	// @Summary Create a new owner
	// @Description Create a new admin user in the system (requires API key)
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param X-API-Key header string true "API key"
	// @Param user body models.UserModel true "Admin user data"
	// @Success 201 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Invalid input"}
	// @Failure 401 {object} gin.H{"error": "Invalid API key"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /admin/create [post]
	router.POST("/users/owner/create", auth.APIKeyMiddleware(), userController.CreateOwner)

	// Prospect APIs
	// @Summary Create a new prospect
	// @Description Create a new prospect in the system
	// @Tags Prospects
	// @Accept json
	// @Produce json
	// @Param prospect body models.ProspectModel true "Prospect data"
	// @Success 201 {object} gin.H{"message": "Prospect created"}
	// @Failure 400 {object} gin.H{"error": "Bad Request"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /prospects [post]
	router.POST("/prospects", prospectController.CreateProspect)

	// @Summary Get a prospect by ID
	// @Description Retrieve a prospect by their unique ID
	// @Tags Prospects
	// @Accept json
	// @Produce json
	// @Param id path string true "Prospect ID"
	// @Success 200 {object} models.ProspectModel
	// @Failure 400 {object} gin.H{"error": "Invalid prospect ID"}
	// @Failure 404 {object} gin.H{"error": "Prospect not found"}
	// @Router /prospects/{id} [get]
	router.GET("/prospects/:id", prospectController.GetProspect)

	// Start the server
	if err := router.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
