package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// MiddlewareManager manages all application middlewares
type MiddlewareManager struct {
	cors      *CORSMiddleware
	csrf      *CSRFMiddleware
	auth      *AuthMiddleware
	rateLimit *RateLimitMiddleware
}

// MiddlewareConfig holds configuration for all middlewares
type MiddlewareConfig struct {
	// Auth configuration
	SessionCookieName string
	SessionMaxAge     time.Duration
	
	// CSRF configuration
	CSRFConfig CSRFConfig
	
	// Rate limit configuration
	RateLimitConfig RateLimitConfig
	
	// Environment
	Environment string
}

// DefaultMiddlewareConfig returns default middleware configuration
func DefaultMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "session_id"),
		SessionMaxAge:     parseDuration("SESSION_MAX_AGE", 24*time.Hour),
		CSRFConfig:        DefaultCSRFConfig(),
		RateLimitConfig:   DefaultRateLimitConfig(),
		Environment:       getEnv("APP_ENV", "development"),
	}
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager(config ...MiddlewareConfig) *MiddlewareManager {
	cfg := DefaultMiddlewareConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Adjust CSRF config for production
	if cfg.Environment == "production" {
		cfg.CSRFConfig.CookieSecure = true
	}

	return &MiddlewareManager{
		cors:      NewCORSMiddleware(),
		csrf:      NewCSRFMiddleware(cfg.CSRFConfig),
		auth:      NewAuthMiddleware(cfg.SessionCookieName, cfg.SessionMaxAge),
		rateLimit: NewRateLimitMiddleware(cfg.RateLimitConfig),
	}
}

// SetupGlobalMiddleware configures global middlewares that apply to all routes
func (m *MiddlewareManager) SetupGlobalMiddleware(app *fiber.App) {
	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency} - ${ip}\n",
		TimeFormat: "15:04:05",
		Output:     os.Stdout,
	}))

	// Recover middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Compression middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// CORS middleware
	app.Use(m.cors.Handler())

	// Rate limiting middleware
	app.Use(m.rateLimit.Handler())

	// Security headers middleware
	app.Use(m.securityHeaders())

	// Custom middleware for API responses
	app.Use(func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})
}

// SetupCSRF sets up CSRF protection for specific routes
func (m *MiddlewareManager) SetupCSRF() fiber.Handler {
	return m.csrf.Handler()
}

// GetAuthMiddleware returns the authentication middleware
func (m *MiddlewareManager) GetAuthMiddleware() *AuthMiddleware {
	return m.auth
}

// RequireAuth returns authentication middleware handler
func (m *MiddlewareManager) RequireAuth() fiber.Handler {
	return m.auth.RequireAuth()
}

// OptionalAuth returns optional authentication middleware handler
func (m *MiddlewareManager) OptionalAuth() fiber.Handler {
	return m.auth.OptionalAuth()
}

// RequireRole returns role-based authorization middleware handler
func (m *MiddlewareManager) RequireRole(role string) fiber.Handler {
	return m.auth.RequireRole(role)
}

// GetCSRFToken returns the CSRF token for the current request
func (m *MiddlewareManager) GetCSRFToken(c fiber.Ctx) string {
	return m.csrf.GetToken(c)
}

// Login authenticates a user and creates a session
func (m *MiddlewareManager) Login(c fiber.Ctx, userID string, userData map[string]interface{}) error {
	return m.auth.Login(c, userID, userData)
}

// Logout logs out a user and destroys the session
func (m *MiddlewareManager) Logout(c fiber.Ctx) error {
	return m.auth.Logout(c)
}

// securityHeaders adds security headers to responses
func (m *MiddlewareManager) securityHeaders() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")
		
		// Enable XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy (basic)
		c.Set("Content-Security-Policy", "default-src 'self'")
		
		return c.Next()
	}
}

// CleanupMiddlewares performs cleanup operations for middlewares
func (m *MiddlewareManager) CleanupMiddlewares() {
	// Cleanup expired sessions
	m.auth.CleanupExpiredSessions()
	
	// Cleanup expired rate limiters
	m.rateLimit.CleanupExpiredLimiters()
}

// StartCleanupRoutine starts a goroutine to periodically cleanup middlewares
func (m *MiddlewareManager) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			m.CleanupMiddlewares()
		}
	}()
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parseDuration parses duration from environment variable with fallback
func parseDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	
	// Try parsing as duration string (e.g., "24h", "30m")
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	
	// Try parsing as hours
	if hours, err := strconv.Atoi(value); err == nil {
		return time.Duration(hours) * time.Hour
	}
	
	return fallback
}