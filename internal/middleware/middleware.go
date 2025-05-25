package middleware

import (
  "github.com/gofiber/fiber/v3"
)

func LoadMiddleware(app *fiber.App) {
	app.Use(ContentTypeHtmlMiddleware())
}