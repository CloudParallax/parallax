package middleware

import (
  "github.com/gofiber/fiber/v3"
)

func ContentTypeHtmlMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Next()
	}
}