package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/cloudparallax/parallax/internal/services"
	"github.com/gofiber/fiber/v3"
)

var preloadService *services.PreloadService

func init() {
	// Initialize preload service
	isProduction := os.Getenv("ENV") == "production"
	preloadService = services.NewPreloadService("./web/static", isProduction)
}

// AssetPreloadMiddleware adds Link headers for preloading CSS and other assets
func AssetPreloadMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		// Only add preload headers for HTML requests (not for static assets)
		if !strings.HasPrefix(path, "/static/dist/css") && !strings.HasPrefix(path, "/api") {
			// Send Early Hints (HTTP 103) before processing the request
			sendEarlyHints(c)

			// Get preload headers from service
			preloadHeaders := preloadService.GetPreloadHeaders()

			if len(preloadHeaders) > 0 {
				c.Set("Link", strings.Join(preloadHeaders, ", "))
			}
		}

		return c.Next()
	}
}

// HTTP2PushMiddleware provides HTTP/2 server push for critical assets
func HTTP2PushMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		// Only push for HTML requests
		if !strings.HasPrefix(path, "/static") && !strings.HasPrefix(path, "/api") {
			// Check if we can push (HTTP/2 support varies by deployment)
			if pusher, ok := c.Response().BodyWriter().(interface {
				Push(target string, opts *http.PushOptions) error
			}); ok {
				// Get push targets from service
				pushTargets := preloadService.GetPushTargets()

				for _, target := range pushTargets {
					if asset, exists := preloadService.GetAsset(target); exists {
						pusher.Push(target, &http.PushOptions{
							Method: "GET",
							Header: http.Header{
								"Content-Type": []string{asset.ContentType},
							},
						})
					}
				}
			}
		}

		return c.Next()
	}
}

// sendEarlyHints sends HTTP 103 Early Hints for critical assets
func sendEarlyHints(c fiber.Ctx) {
	// Get preload headers for early hints
	preloadHeaders := preloadService.GetPreloadHeaders()

	if len(preloadHeaders) > 0 {
		// Try to send Early Hints (HTTP 103)
		// Note: This requires server support for HTTP 103
		if writer, ok := c.Response().BodyWriter().(http.ResponseWriter); ok {
			writer.Header().Set("Link", strings.Join(preloadHeaders, ", "))
			writer.WriteHeader(103) // HTTP 103 Early Hints
		}
	}
}
