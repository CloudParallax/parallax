package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v3"
)

// CSRFMiddleware provides CSRF protection using double submit cookie pattern
type CSRFMiddleware struct {
	tokenLookup    string
	cookieName     string
	headerName     string
	cookieSecure   bool
	cookieHTTPOnly bool
	cookieSameSite string
	expiration     time.Duration
	keyGenerator   func() ([]byte, error)
}

// CSRFConfig holds CSRF middleware configuration
type CSRFConfig struct {
	TokenLookup    string              // "header:X-CSRF-Token" or "form:_token"
	CookieName     string              // Name of the CSRF cookie
	HeaderName     string              // Name of the CSRF header
	CookieSecure   bool                // Set cookie secure flag
	CookieHTTPOnly bool                // Set cookie httponly flag
	CookieSameSite string              // Set cookie samesite
	Expiration     time.Duration       // Token expiration time
	KeyGenerator   func() ([]byte, error)
}

// DefaultCSRFConfig returns default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		HeaderName:     "X-CSRF-Token",
		CookieSecure:   false, // Should be true in production with HTTPS
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour,
		KeyGenerator:   generateRandomBytes,
	}
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware(config ...CSRFConfig) *CSRFMiddleware {
	cfg := DefaultCSRFConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return &CSRFMiddleware{
		tokenLookup:    cfg.TokenLookup,
		cookieName:     cfg.CookieName,
		headerName:     cfg.HeaderName,
		cookieSecure:   cfg.CookieSecure,
		cookieHTTPOnly: cfg.CookieHTTPOnly,
		cookieSameSite: cfg.CookieSameSite,
		expiration:     cfg.Expiration,
		keyGenerator:   cfg.KeyGenerator,
	}
}

// Handler returns the CSRF middleware handler
func (c *CSRFMiddleware) Handler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// Skip CSRF for safe methods
		if c.isSafeMethod(ctx.Method()) {
			return c.setCSRFToken(ctx)
		}

		// Validate CSRF token for unsafe methods
		if !c.validateCSRFToken(ctx) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusForbidden,
					"message": "CSRF token mismatch",
				},
			})
		}

		return ctx.Next()
	}
}

// isSafeMethod checks if HTTP method is safe (doesn't modify state)
func (c *CSRFMiddleware) isSafeMethod(method string) bool {
	return method == fiber.MethodGet ||
		method == fiber.MethodHead ||
		method == fiber.MethodOptions
}

// setCSRFToken generates and sets a new CSRF token
func (c *CSRFMiddleware) setCSRFToken(ctx fiber.Ctx) error {
	// Check if token already exists and is valid
	existingToken := ctx.Cookies(c.cookieName)
	if existingToken != "" && c.isValidToken(existingToken) {
		// Token exists and is valid, add to response headers for client access
		ctx.Set(c.headerName, existingToken)
		return ctx.Next()
	}

	// Generate new token
	token, err := c.generateToken()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to generate CSRF token",
			},
		})
	}

	// Set CSRF cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     c.cookieName,
		Value:    token,
		Expires:  time.Now().Add(c.expiration),
		HTTPOnly: c.cookieHTTPOnly,
		Secure:   c.cookieSecure,
		SameSite: c.cookieSameSite,
	})

	// Add token to response headers for client access
	ctx.Set(c.headerName, token)

	return ctx.Next()
}

// validateCSRFToken validates the CSRF token using double submit cookie pattern
func (c *CSRFMiddleware) validateCSRFToken(ctx fiber.Ctx) bool {
	// Get token from cookie
	cookieToken := ctx.Cookies(c.cookieName)
	if cookieToken == "" {
		return false
	}

	// Get token from header or form
	var headerToken string
	if c.tokenLookup == "header:"+c.headerName {
		headerToken = ctx.Get(c.headerName)
	} else {
		// Parse form token if needed
		headerToken = ctx.FormValue("_token")
	}

	if headerToken == "" {
		return false
	}

	// Validate both tokens exist and match
	if !c.isValidToken(cookieToken) || !c.isValidToken(headerToken) {
		return false
	}

	// Use constant time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) == 1
}

// generateToken generates a new CSRF token
func (c *CSRFMiddleware) generateToken() (string, error) {
	bytes, err := c.keyGenerator()
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// isValidToken checks if token format is valid
func (c *CSRFMiddleware) isValidToken(token string) bool {
	if token == "" {
		return false
	}
	
	// Decode token to validate format
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	
	// Check minimum length (32 bytes = 256 bits)
	return len(decoded) >= 32
}

// generateRandomBytes generates cryptographically secure random bytes
func generateRandomBytes() ([]byte, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	return bytes, err
}

// GetToken returns the current CSRF token for the request
func (c *CSRFMiddleware) GetToken(ctx fiber.Ctx) string {
	return ctx.Cookies(c.cookieName)
}

// InvalidateToken removes the CSRF token cookie
func (c *CSRFMiddleware) InvalidateToken(ctx fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     c.cookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: c.cookieHTTPOnly,
		Secure:   c.cookieSecure,
		SameSite: c.cookieSameSite,
	})
}