package handlers

import (
	"fmt"
	"strings"

	"github.com/cloudparallax/parallax/internal/services"
	"github.com/cloudparallax/parallax/internal/views"
	"github.com/cloudparallax/parallax/web/templates"
	"github.com/cloudparallax/parallax/web/templates/layouts"
	"github.com/cloudparallax/parallax/web/templates/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

var blogService *services.BlogService

func LoadHandlers(app *fiber.App) {
	// Initialize blog service
	blogService = services.NewBlogService("content/blog")

	// Helper function to check if request is from HTMX
	isHTMXRequest := func(c fiber.Ctx) bool {
		return c.Get("HX-Request") == "true"
	}

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
	// app.Get("/contact", func(c fiber.Ctx) error {
	// 	csrfToken := csrf.TokenFromContext(c)
	// 	component := pages.ContactPage()

	// 	if isHTMXRequest(c) {
	// 		return component.Render(c.Context(), c.Response().BodyWriter())
	// 	}

	// 	return layouts.BaseLayout("Contact", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	// })

	// Demo API routes for HTMX examples
	app.Get("/api/demo/content/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		var content string
		switch id {
		case "1":
			content = `
				<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
					<h3 class="text-lg font-semibold text-blue-900 mb-2">Dynamic Content 1</h3>
					<p class="text-blue-700">This content was loaded dynamically using HTMX!
					   Notice how the page didn't refresh, but the content updated seamlessly.</p>
					<div class="mt-3">
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
							HTMX Powered
						</span>
					</div>
				</div>
			`
		case "2":
			content = `
				<div class="bg-green-50 border border-green-200 rounded-lg p-4">
					<h3 class="text-lg font-semibold text-green-900 mb-2">Dynamic Content 2</h3>
					<p class="text-green-700">This is a different piece of content, also loaded via HTMX.
					   The server can return any HTML content, making it very flexible.</p>
					<ul class="mt-3 list-disc list-inside text-green-600 text-sm">
						<li>No page refresh required</li>
						<li>Server-rendered content</li>
						<li>Progressive enhancement</li>
					</ul>
				</div>
			`
		case "3":
			content = `
				<div class="bg-purple-50 border border-purple-200 rounded-lg p-4">
					<h3 class="text-lg font-semibold text-purple-900 mb-2">Dynamic Content 3</h3>
					<p class="text-purple-700">And here's yet another example! HTMX makes it easy to create
					   interactive user interfaces without complex JavaScript frameworks.</p>
					<div class="mt-3 flex space-x-2">
						<button class="bg-purple-600 hover:bg-purple-700 text-white text-xs font-medium py-1 px-2 rounded"
						        hx-get="/api/demo/content/1" hx-target="#demo-content">
							Load Content 1
						</button>
						<button class="bg-purple-600 hover:bg-purple-700 text-white text-xs font-medium py-1 px-2 rounded"
						        hx-get="/api/demo/content/2" hx-target="#demo-content">
							Load Content 2
						</button>
					</div>
				</div>
			`
		default:
			content = `
				<div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
					<p class="text-gray-500">Unknown content ID</p>
				</div>
			`
		}

		// Clean up extra whitespace
		content = strings.TrimSpace(content)
		return c.Type("text/html").SendString(content)
	})

	// Legacy routes for backward compatibility
	app.Get("/test2", func(c fiber.Ctx) error {
		// Get CSRF token from the context
		token := csrf.TokenFromContext(c)
		if token == "" {
			panic("CSRF token not found in context")
		}

		// Note: Make sure this matches the KeyLookup configured in your middleware
		formKey := "_csrf"

		// Create a form with the CSRF token
		tmpl := fmt.Sprintf(`<form action="/post" method="POST">
        <input type="hidden" name="%s" value="%s">
        <input type="text" name="message">
        <input type="submit" value="Submit">
    </form>`, formKey, token)

		c.Set("Content-Type", "text/html")
		return c.SendString(tmpl)
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
