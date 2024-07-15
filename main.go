package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/riwayatt1/golang-mygram/config"
	"github.com/riwayatt1/golang-mygram/handlers"
	"github.com/riwayatt1/golang-mygram/middleware"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()

	// Initialize database connection
	config.ConnectDatabase()

	r.POST("/user/register", handlers.RegisterUser)
	r.POST("/user/login", handlers.LoginUser)

	// Routes
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())

	{
		// User endpoints
		api.GET("/user/profile", middleware.AuthorizeUser(), handlers.GetUserProfile)
		api.PUT("/users/:userId", middleware.AuthorizeUser(), handlers.UpdateUser)
		api.DELETE("/users/:userId", middleware.AuthorizeUser(), handlers.DeleteUser)

		// Photo endpoints
		api.POST("/photo/create", handlers.CreatePhoto)
		api.GET("/photo", handlers.GetAllPhotos)
		api.GET("/photo/:id", middleware.AuthorizeUser(), handlers.GetPhoto)
		api.PUT("/photo/:id", handlers.UpdatePhoto)
		api.DELETE("/photo/:id", handlers.DeletePhoto)

		// Comment endpoints
		api.POST("/comments/create", middleware.AuthorizeUser(), handlers.CreateComment)
		api.GET("/comments/:id", middleware.AuthorizeUser(), handlers.GetComment)
		api.PUT("/comments/:id", middleware.AuthorizeUser(), handlers.UpdateComment)
		api.DELETE("/comments/:id", middleware.AuthorizeUser(), handlers.DeleteComment)

		// SocialMedia endpoints
		api.POST("/socialmedias/create", middleware.AuthorizeUser(), handlers.CreateSocialMedia)
		api.GET("/socialmedias/:id", middleware.AuthorizeUser(), handlers.GetSocialMedia)
		api.PUT("/socialmedias/:id", middleware.AuthorizeUser(), handlers.UpdateSocialMedia)
		api.DELETE("/socialmedias/:id", middleware.AuthorizeUser(), handlers.DeleteSocialMedia)
	}

	// Run the server
	r.Run(":8080")
}
