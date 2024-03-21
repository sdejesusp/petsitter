package main

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Swagger configuration
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/petsitter.yaml",
		Path:     "swagger",
		Title:    "Petsitter API Docs",
		CacheAge: 1,
	}

	app.Use(swagger.New(cfg))

	app.Get("/welcome", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the PetSitter API")
	})

	app.Listen(":3000")
}
