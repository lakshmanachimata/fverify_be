package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"fverify_be/internal/auth"
	"fverify_be/internal/controllers"
	"fverify_be/internal/repositories"
	"fverify_be/internal/services"

	"fverify_be/cmd/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// @title FVerify API
// @version 1.0
// @description This is the API documentation for the Fverify backend.
// @host localhost:9000
// @BasePath /
// @schemes http
func main() {
	// Load configuration
	viper.SetConfigName("config_db") // Name of the config file (without extension)
	viper.SetConfigType("json")      // Config file type
	viper.AddConfigPath(".")         // Path to look for the config file in the current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get MongoDB credentials from config
	username := viper.GetString("mongodb.username")
	password := viper.GetString("mongodb.password")
	uri := viper.GetString("mongodb.uri")
	// URL-encode the username and password
	encodedUsername := url.QueryEscape(username)
	encodedPassword := url.QueryEscape(password)

	// Construct MongoDB URI
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s", encodedUsername, encodedPassword, uri)

	// Set up MongoDB connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

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
	prospectRepo := repositories.NewProspectRepository(client, "fverify_db", "prospects")
	userRepo := repositories.NewUserRepository(client, "fverify_db", "users")
	orgRepo := repositories.NewOrganisationRepository(client, "fverify_db", "orgs")

	// Initialize services
	prospectService := services.NewProspectService(prospectRepo)
	userService := services.NewUserService(userRepo)
	orgService := services.NewOrganisationService(orgRepo, userRepo)

	// Initialize controllers
	prospectController := controllers.NewProspectController(prospectService)
	userController := controllers.NewUserController(userService, orgService)
	organisationController := controllers.NewOrganisationController(orgService)

	// Set up Gin router
	router := gin.Default()
	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key", "org_id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger setup
	docs.SwaggerInfo.Title = "FVerify API"
	docs.SwaggerInfo.Description = "This is the API documentation for the Fverify backend."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:9000"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/organisations", auth.OrgAPIKeyMiddleware(), organisationController.CreateOrganisation)
		api.PUT("/organisations/:org_id", auth.OrgAPIKeyMiddleware(), organisationController.UpdateOrganisation)
		// api.DELETE("/organisations/:org_id", auth.OrgAPIKeyMiddleware(), organisationController.DeleteOrganisation)
		api.GET("/organisations", auth.OrgAPIKeyMiddleware(), organisationController.GetAllOrganisations)
		api.POST("/users/login", userController.LoginUser)
		api.POST("/users", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead"), userController.CreateUser)
		api.PUT("/users/uid/:uId", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive"), userController.UpdateUser)
		api.GET("/users", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive"), userController.GetAllUsers)
		api.GET("/users/:userId", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive"), userController.GetUserByUserID)
		// api.DELETE("/users/uid/:uId", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner"), userController.DeleteUserByUId)
		// api.DELETE("/users/userid/:userId", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner"), userController.DeleteUserByUserId)
		// api.PUT("/users/uid/:uId/setpassword", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive"), userController.SetPassword)
		api.POST("/users/admin/create", auth.APIKeyMiddleware(), userController.CreateAdmin)
		api.POST("/users/owner/create", auth.APIKeyMiddleware(), userController.CreateOwner)
		api.GET("/users/roles", userController.GetUserRoles)
		api.GET("/users/statuses", userController.GetUserStatuses)
		api.POST("/prospects", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive", "Field Lead"), prospectController.CreateProspect)
		api.GET("/prospects/:uid", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive", "Field Lead", "Field Executive"), prospectController.GetProspect)
		api.PUT("/prospects", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive", "Field Lead", "Field Executive"), prospectController.UpdateProspect)
		api.GET("/prospects", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive", "Field Lead", "Field Executive"), prospectController.GetProspects)
		api.GET("/prospects/count", auth.AuthMiddleware(*orgRepo, *userRepo, "Admin", "Owner", "Operations Lead", "Operations Executive", "Field Lead", "Field Executive"), prospectController.GetProspectsCount)
	}

	// Start the server
	if err := router.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
