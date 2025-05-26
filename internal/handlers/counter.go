package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

var counterValue int = 1

// loadCounterHandler sets up counter-related routes
func loadCounterHandler(app *fiber.App) {
	counterGroup := app.Group("/counter")

	// Increment counter endpoint
	counterGroup.Put("/increment", handleCounterIncrement)
	
	// Decrement counter endpoint  
	counterGroup.Put("/decrement", handleCounterDecrement)
	
	// Get current counter value
	counterGroup.Get("/", handleCounterGet)
}

// handleCounterIncrement increments the counter value
func handleCounterIncrement(c fiber.Ctx) error {
	counterValue++
	return c.SendString(fmt.Sprintf("%d", counterValue))
}

// handleCounterDecrement decrements the counter value (minimum value is 1)
func handleCounterDecrement(c fiber.Ctx) error {
	if counterValue > 1 {
		counterValue--
	}
	return c.SendString(fmt.Sprintf("%d", counterValue))
}

// handleCounterGet returns the current counter value
func handleCounterGet(c fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("%d", counterValue))
}