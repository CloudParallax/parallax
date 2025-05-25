package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func CSRFMiddleware(sessionStore *session.Store) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:         "header:X-Csrf-Token",
		CookieName:        "csrf_",
		CookieSameSite:    "Lax",
		CookieSecure:      false, // Set to true in production with HTTPS
		CookieHTTPOnly:    true,
		CookieSessionOnly: true,
		Session:           sessionStore,
		IdleTimeout:       30 * time.Minute,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Check if request is HTMX request
			if c.Get("HX-Request") == "true" {
				// Return HTMX-friendly error response
				c.Set("HX-Retarget", "body")
				c.Set("HX-Reswap", "innerHTML")
				return c.Status(fiber.StatusForbidden).SendString(`
					<div class="bg-red-50 border border-red-200 rounded-lg p-4 m-4">
						<h3 class="text-lg font-medium text-red-800 mb-2">CSRF Token Error</h3>
						<p class="text-red-700">Your session has expired. Please refresh the page and try again.</p>
						<button onclick="window.location.reload()"
							class="mt-3 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700">
							Refresh Page
						</button>
					</div>
				`)
			}

			// For regular requests, redirect to error page or return JSON
			if c.Get("Content-Type") == "application/json" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":   "CSRF token validation failed",
					"code":    "CSRF_ERROR",
					"message": "Please refresh the page and try again",
				})
			}

			// Default error response
			return c.Status(fiber.StatusForbidden).SendString("CSRF token validation failed. Please refresh the page and try again.")
		},
	})
}
