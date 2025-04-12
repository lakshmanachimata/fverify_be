package main

import (
	"context"
	"log"

	"kowtha_be/internal/controllers"
	"kowtha_be/internal/repositories"
	"kowtha_be/internal/services"

	"kowtha_be/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Kowtha API
// @version 1.0
// @description This is the API documentation for the Kowtha backend.
// @host localhost:9000
// @BasePath /
// @schemes http
func main() {
	// Set up MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize repositories
	prospectRepo := repositories.NewProspectRepository(client, "kowtha_db", "prospects")
	userRepo := repositories.NewUserRepository(client, "kowtha_db", "users")

	// Initialize services
	prospectService := services.NewProspectService(prospectRepo)
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	prospectController := controllers.NewProspectController(prospectService)
	userController := controllers.NewUserController(userService)

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

	// User APIs
	// @Summary Create a new user
	// @Description Create a new user in the system
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param user body models.UserModel true "User data"
	// @Success 201 {object} gin.H{"message": "User created successfully"}
	// @Failure 400 {object} gin.H{"error": "Bad Request"}
	// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
	// @Router /users [post]
	router.POST("/users", userController.CreateUser)

	// @Summary Get a user by ID
	// @Description Retrieve a user by their unique ID
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Param id path string true "User ID"
	// @Success 200 {object} models.UserModel
	// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Router /users/{id} [get]
	router.GET("/users/:id", userController.GetUserByID)

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
