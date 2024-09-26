package models

import (
	"fmt" // Importing fmt package for formatted I/O operations
	"log"  // Importing log package for logging errors

	"gorm.io/driver/postgres" // Importing PostgreSQL driver for GORM
	"gorm.io/gorm"            // Importing GORM for ORM functionalities
)

var (
	DB  *gorm.DB // Global variable to hold the database connection
	err error     // Global variable to hold errors
)

// ConnectDatabase establishes a connection to the PostgreSQL database.
// Parameters:
// - host: The database server host address
// - user: The username for the database
// - password: The password for the database user
// - dbName: The name of the database to connect to
// - port: The port number for the database server
// - sslmode: The SSL mode for the connection (e.g., "disable", "require")
func ConnectDatabase(host, user, password, dbName string, port int, sslmode string) {
	// Create a Data Source Name (DSN) string to configure the database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbName, port, sslmode)
	
	// Open a new database connection using GORM with the DSN
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Log a fatal error if the connection fails
		log.Fatal("Error connecting to database:", err)
	}

	// Auto migrate the schema for the File model
	// This will create the table if it does not exist and update the table structure if necessary
	err = DB.AutoMigrate(&File{})
	if err != nil {
		// Log a fatal error if the migration fails
		log.Fatal("Error migrating schema:", err)
	}
}
