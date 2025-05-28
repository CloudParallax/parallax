package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// CORSMiddleware provides CORS functionality
type CORSMiddleware struct {
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           int
}

// NewCORSMiddleware creates a new CORS middleware with environment configuration
func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{
		allowOrigins:     parseEnvArray("CORS_ALLOW_ORIGINS", []string{"*"}),
		allowMethods:     parseEnvArray("CORS_ALLOW_METHODS", []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		allowHeaders:     parseEnvArray("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}),
		exposeHeaders:    parseEnvArray("CORS_EXPOSE_HEADERS", []string{"Content-Length"}),
		allowCredentials: parseEnvBool("CORS_ALLOW_CREDENTIALS", false),
		maxAge:           parseEnvInt("CORS_MAX_AGE", 86400),
	}
}

// Handler returns the CORS middleware handler
func (c *CORSMiddleware) Handler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		origin := ctx.Get("Origin")
		
		// Handle preflight requests
		if ctx.Method() == fiber.MethodOptions {
			return c.handlePreflight(ctx, origin)
		}

		// Handle actual requests
		return c.handleActualRequest(ctx, origin)
	}
}

// handlePreflight handles CORS preflight requests
func (c *CORSMiddleware) handlePreflight(ctx fiber.Ctx, origin string) error {
	// Check if origin is allowed
	if !c.isOriginAllowed(origin) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusForbidden,
				"message": "Origin not allowed",
			},
		})
	}

	// Set CORS headers for preflight
	c.setAllowOriginHeader(ctx, origin)
	ctx.Set("Access-Control-Allow-Methods", strings.Join(c.allowMethods, ", "))
	ctx.Set("Access-Control-Allow-Headers", strings.Join(c.allowHeaders, ", "))
	
	if c.allowCredentials {
		ctx.Set("Access-Control-Allow-Credentials", "true")
	}
	
	if c.maxAge > 0 {
		ctx.Set("Access-Control-Max-Age", strconv.Itoa(c.maxAge))
	}

	// Respond to preflight with 204 No Content
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}

// handleActualRequest handles actual CORS requests
func (c *CORSMiddleware) handleActualRequest(ctx fiber.Ctx, origin string) error {
	// Check if origin is allowed
	if origin != "" && !c.isOriginAllowed(origin) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusForbidden,
				"message": "Origin not allowed",
			},
		})
	}

	// Set CORS headers for actual request
	if origin != "" {
		c.setAllowOriginHeader(ctx, origin)
	}
	
	if len(c.exposeHeaders) > 0 {
		ctx.Set("Access-Control-Expose-Headers", strings.Join(c.exposeHeaders, ", "))
	}
	
	if c.allowCredentials {
		ctx.Set("Access-Control-Allow-Credentials", "true")
	}

	return ctx.Next()
}

// isOriginAllowed checks if the origin is in the allowed list
func (c *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return true // Allow requests without Origin header
	}

	for _, allowed := range c.allowOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
		// Support wildcard subdomains like *.example.com
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*.")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}

// setAllowOriginHeader sets the appropriate Access-Control-Allow-Origin header
func (c *CORSMiddleware) setAllowOriginHeader(ctx fiber.Ctx, origin string) {
	if c.allowCredentials {
		// When credentials are allowed, we can't use "*"
		// We must echo back the specific origin
		if origin != "" && c.isOriginAllowed(origin) {
			ctx.Set("Access-Control-Allow-Origin", origin)
		}
	} else {
		// When credentials are not allowed, we can use "*" or specific origin
		for _, allowed := range c.allowOrigins {
			if allowed == "*" {
				ctx.Set("Access-Control-Allow-Origin", "*")
				return
			}
		}
		if origin != "" && c.isOriginAllowed(origin) {
			ctx.Set("Access-Control-Allow-Origin", origin)
		}
	}
}

// parseEnvArray parses comma-separated environment variable into string slice
func parseEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parts := strings.Split(value, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	
	return result
}

// parseEnvBool parses environment variable as boolean
func parseEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}

// parseEnvInt parses environment variable as integer
func parseEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}