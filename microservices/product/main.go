package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		return nil, fmt.Errorf("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		return db, fmt.Errorf("error connecting to the database: %v", err)
	}

	InitModels(db)

	return db, nil
}

func main() {

	// --------------------------------------------------------------------------
	// Load environment variables
	//
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading [product-service] .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Invalid port - not found in environment")
	}

	// --------------------------------------------------------------------------
	// Init database
	//
	_, err := initDB()

	if err != nil {
		log.Fatalf("Error connecting to [product-service] database: %v", err)
	} else {
		fmt.Println("*** >>> Successfully initialized [Postgres]")
	}

	// --------------------------------------------------------------------------
	// Init Fiber
	//
	app := fiber.New()

	// --------------------------------------------------------------------------
	// Routes
	//
	app.Post("/create", createProduct)
	app.Post("/delete", deleteProduct)

	// --------------------------------------------------------------------------
	// Launch user service
	//
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
