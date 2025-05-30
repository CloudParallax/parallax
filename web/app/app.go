package app

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudparallax/parallax/internal/adapters/http"
	"github.com/gofiber/fiber/v3"
)

// LoadApp initializes and starts the API server
func LoadApp() {
	port := GetEnv("SERVER_PORT", "8080")
	log.Printf("Starting Parallax API server on port %s\n", port)

	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		AppName:      "Parallax API v1.0.0",
	})

	// Setup routes (router handles all middleware setup)
	router := http.NewRouter(app)
	log.Println("Setting up routes...")
	router.SetupRoutes()
	log.Println("Routes setup complete")

	fmt.Printf("🚀 Starting Parallax API server at :%s\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

// customErrorHandler handles errors in a consistent way
func customErrorHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Printf("Error occurred: %s (code: %d)\n", err.Error(), code)
	return ctx.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}

// GetEnv gets environment variable with fallback
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Environment variable %s is set to %s\n", key, value)
		return value
	}
	log.Printf("Environment variable %s is not set, using fallback %s\n", key, fallback)
	return fallback
}