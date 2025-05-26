package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func LoadMiddleware(app *fiber.App) {
	sessionMiddleware, sessionStore := session.NewWithStore()
	app.Use(CorsMiddleware())
	app.Use(AssetPreloadMiddleware())
	app.Use(sessionMiddleware)
	app.Use(CSRFMiddleware(sessionStore))
	app.Use(ContentTypeHtmlMiddleware())
}
