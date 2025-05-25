package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CorsMiddleware() fiber.Handler {
	// Get CORS configuration from environment variables with defaults
	allowOrigins := getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"*"})
	allowMethods := getEnvAsSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"})
	allowHeaders := getEnvAsSlice("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"})
	exposeHeaders := getEnvAsSlice("CORS_EXPOSE_HEADERS", []string{"Content-Length"})
	allowCredentials := getEnvAsBool("CORS_ALLOW_CREDENTIALS", false)
	maxAge := getEnvAsInt("CORS_MAX_AGE", 86400) // 24 hours default

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: allowCredentials,
		ExposeHeaders:    exposeHeaders,
		MaxAge:           maxAge,
	})
}

// Helper function to get environment variable as string slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// Helper function to get environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}