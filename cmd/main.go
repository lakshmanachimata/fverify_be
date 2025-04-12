package main

import (
	"context"
	"log"

	"kowtha_be/internal/controllers"
	"kowtha_be/internal/repositories"
	"kowtha_be/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	router.POST("/users", userController.CreateUser)
	router.GET("/users/:id", userController.GetUserByID)

	router.POST("/prospects", prospectController.CreateProspect)
	router.GET("/prospects/:id", prospectController.GetProspect)

	// Start the server
	if err := router.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
