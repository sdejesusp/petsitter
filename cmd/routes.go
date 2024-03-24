package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/petsitter/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/welcome", handlers.WelcomeMessage)

	// user routes
	users := app.Group("/users")
	users.Post("/", handlers.CreateUser)
	users.Get("/:id", handlers.GetUserWithId)
}
