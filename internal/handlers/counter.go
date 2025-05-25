package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

var counterValue int = 1

func loadCounterHandler(app *fiber.App) {
	counterGroup := app.Group("/counter")

	counterGroup.Put("/increment", func(c fiber.Ctx) error {
		counterValue++
		return c.SendString(fmt.Sprintf("%d", counterValue))
	})

	counterGroup.Put("/decrement", func(c fiber.Ctx) error {
		if counterValue != 1 {
			counterValue--
		}
		return c.SendString(fmt.Sprintf("%d", counterValue))
	})
}