package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/petsitter/handlers"
)

func setupRoutesWithoutJWT(app *fiber.App) {
	app.Get("/welcome", handlers.WelcomeMessage)
	app.Post("/token", handlers.GetToken)
	users := app.Group("/users")
	users.Post("/", handlers.CreateUser)
}

func setupRoutes(app *fiber.App) {
	// user routes
	users := app.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/profile", handlers.GetUserProfile)
	users.Get("/:id", handlers.GetUserWithId)
	users.Put("/:id", handlers.ModifyUserWithId)
	users.Delete("/:id", handlers.DeleteUserWithId)
}
