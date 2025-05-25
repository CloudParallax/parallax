package handlers

import (
	"github.com/cloudparallax/parallax/internal/views"
	"github.com/cloudparallax/parallax/web/templates"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

func LoadHandlers(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := templates.Splash()
		return templates.Layout(component, csrfToken).Render(c.Context(), c.Response().BodyWriter())
	})

	app.Get("/test", func(c fiber.Ctx) error {
		return views.Render(c, templates.Splash())
	})

	// Test form handler for CSRF demonstration
	app.Post("/api/test-form", func(c fiber.Ctx) error {
		message := c.FormValue("message")
		if message == "" {
			message = "No message provided"
		}
		
		response := `<div class="text-green-700 bg-green-100 border border-green-200 rounded p-2">
			<strong>✅ Form submitted successfully!</strong><br>
			Message: ` + message + `<br>
			<small class="text-green-600">CSRF token was validated automatically</small>
		</div>`
		
		return c.Type("text/html").SendString(response)
	})

	loadCounterHandler(app)
}
