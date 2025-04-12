package main

import (
	// Import the generated docs package
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
	"go.mongodb.org/mongo-driver/mongo/options" // Import the generated docs package
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

	// Initialize repositories and services
	prospectController := controllers.NewProspectController(prospectService)

	// Initialize repositories and services
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

	router.POST("/users", userController.CreateUser)
	router.GET("/users/:id", userController.GetUserByID)

	router.POST("/prospects", prospectController.CreateProspect)
	router.GET("/prospects/:id", prospectController.GetProspect)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	if err := router.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
