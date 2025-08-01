package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go001/cmd/server"        // Your API handler implementation
	"go001/internal/products" // Your business logic service, models, and repository
	"go001/pkg/api"           // Generated API types and interfaces
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen  --config pkg/api/openapichiserver.yaml pkg/api/openapimodels.yaml

func main() {

	// Load environment variables from .env file
	readEnv()

	// Initialize logging
	gormLogger := iniLog()

	// Database Connection (PostgreSQL with GORM) ---
	productService := iniDb(gormLogger)

	// Create an instance of your ServerInterface implementation ---
	iniSrv(productService)
}

func readEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func iniLog() logger.Interface {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable colors, use default
		},
	)
	return gormLogger
}

func iniDb(gormLogger logger.Interface) *products.Service {
	// Get environment variables for database connection
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	log.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// --- 2. Auto-Migrate Database Schema ---
	// In a production environment, you'd use a dedicated migration tool (e.g., golang-migrate).
	// For development, GORM's AutoMigrate is convenient.
	err = db.AutoMigrate(&products.ProductDB{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database migration complete.")

	// Optional: Seed some initial data if the table is empty (for testing)
	var count int64
	db.Model(&products.ProductDB{}).Count(&count)
	if count == 0 {
		log.Println("Seeding initial product data...")
		seedProducts := []products.ProductDB{
			{Name: "Laptop Pro X", Description: "High performance laptop with 16GB RAM and 512GB SSD.", Price: 1200.50, Category: "Electronics"},
			{Name: "Mechanical Keyboard", Description: "Clicky keys for the best typing experience.", Price: 99.99, Category: "Electronics"},
			{Name: "Ergonomic Chair", Description: "Comfortable chair for long working hours.", Price: 350.00, Category: "Furniture"},
		}
		if err := db.Create(&seedProducts).Error; err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}
		log.Println("Initial product data seeded.")
	}

	// --- 3. Initialize Repositories and Services ---
	productRepo := products.NewGORMProductRepository(db)
	productService := products.NewService(productRepo)

	return productService
}

func iniSrv(productService *products.Service) {
	// --- 4. Create an instance of your API handler ---
	apiHandlerImpl := server.NewAPIHandler(productService)

	// --- 5. Create a new Chi router ---
	r := chi.NewRouter()

	// --- 6. Add Chi middleware ---
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- 7. Register the generated handlers with your Chi router ---
	api.HandlerFromMux(apiHandlerImpl, r)

	// --- 8. Create and start the HTTP server ---
	port := ":" + os.Getenv("APP_PORT")
	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped.")
}
