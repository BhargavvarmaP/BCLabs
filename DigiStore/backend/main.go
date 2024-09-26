package main

import (
	"Distributed_File_Storage/controllers" // Importing the controllers for handling requests
	"Distributed_File_Storage/models"      // Importing the models for database interaction
	"github.com/gin-contrib/cors"          // Importing CORS middleware for handling cross-origin requests
	"github.com/gin-gonic/gin"             // Importing the Gin framework for building the web server
	"log"                                   // Importing log package for logging errors
	"time"                                  // Importing time package for handling time-related functionality
)

func main() {
	// Initialize the database connection using the ConnectDatabase function from the models package
	models.ConnectDatabase("localhost", "postgres", "admin", "postgres", 5432, "disable")

	// Set up the Gin router with default middleware
	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Specify allowed origins for CORS
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Specify allowed HTTP methods
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Specify allowed headers in requests
		ExposeHeaders:    []string{"Content-Length"}, // Specify headers to expose to clients
		AllowCredentials: true, // Allow credentials in CORS requests
		MaxAge:           12 * time.Hour, // Cache preflight request results for 12 hours
	}

	// Apply the CORS middleware to the router
	r.Use(cors.New(corsConfig))

	// Define routes and associate them with their respective controller functions
	r.POST("/upload", controllers.UploadFile)    // Route for uploading files
	r.GET("/files", controllers.GetFiles)         // Route for retrieving a list of files
	r.GET("/download/:id", controllers.DownloadFile) // Route for downloading a specific file by ID

	// Run the server on port 8080 and log an error if it fails to start
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err) // Log the error and exit
	}
}
