package handlers

import (
	"os"
	"os/exec"

	"github.com/cloudparallax/parallax/internal/services"
	"github.com/cloudparallax/parallax/web/templates/components"
	"github.com/cloudparallax/parallax/web/templates/layouts"
	"github.com/cloudparallax/parallax/web/templates/pages"
	"github.com/cloudparallax/parallax/web/templates/responses"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

var blogService *services.BlogService

// LoadHandlers sets up all application routes and handlers
func LoadHandlers(app *fiber.App) {
	// Initialize services
	blogService = services.NewBlogService("content/blog")

	// Register route groups
	setupMainRoutes(app)
	setupAPIRoutes(app)
	setupDemoRoutes(app)
	
	// Load additional handlers
	loadCounterHandler(app)
}

// setupMainRoutes configures the main application routes
func setupMainRoutes(app *fiber.App) {
	// Home page
	app.Get("/", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.HomePage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Home", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// Blog routes
	app.Get("/blog", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		posts, err := blogService.GetAllPosts()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error loading blog posts")
		}

		component := pages.BlogListPage(posts)

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Blog", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	app.Get("/blog/:slug", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		slug := c.Params("slug")

		post, err := blogService.GetPostBySlug(slug)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Blog post not found")
		}

		component := pages.BlogPostPage(*post)

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout(post.Title, csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// About page
	app.Get("/about", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.AboutPage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("About", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// Contact page
	app.Get("/contact", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.ContactPage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Contact", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})
}

// setupAPIRoutes configures API endpoints
func setupAPIRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Form handler for CSRF demonstration
	api.Post("/test-form", func(c fiber.Ctx) error {
		message := c.FormValue("message")
		if message == "" {
			message = "No message provided"
		}

		component := responses.FormSuccessResponse(message)
		return component.Render(c.Context(), c.Response().BodyWriter())
	})

	// CSS rebuild endpoint for development
	api.Post("/rebuild-css", func(c fiber.Ctx) error {
		if os.Getenv("ENV") == "production" {
			return c.Status(fiber.StatusForbidden).SendString("Not available in production")
		}

		cmd := exec.Command("tailwindcss", "-i", "./web/static/main.css", "-o", "./web/static/dist/output.css")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			errorMsg := `<div class="text-red-700 bg-red-100 border border-red-200 rounded p-2">
				<strong>❌ CSS rebuild failed!</strong><br/>
				Error: ` + err.Error() + `<br/>
				<pre class="text-xs mt-2 text-red-600">` + string(output) + `</pre>
			</div>`
			return c.Status(fiber.StatusInternalServerError).Type("text/html").SendString(errorMsg)
		}

		component := responses.BuildSuccessResponse(string(output))
		return component.Render(c.Context(), c.Response().BodyWriter())
	})
}

// setupDemoRoutes configures demo routes for HTMX examples
func setupDemoRoutes(app *fiber.App) {
	demo := app.Group("/api/demo")

	// Dynamic content for HTMX demonstrations
	demo.Get("/content/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		component := components.DemoContent(id)
		
		return component.Render(c.Context(), c.Response().BodyWriter())
	})
}

// isHTMXRequest checks if the request comes from HTMX
func isHTMXRequest(c fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}