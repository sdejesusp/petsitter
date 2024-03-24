package main

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/petsitter/database"
)

func main() {
	// Stablish connection to the DB
	database.ConnectDb()

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

	// API routes
	setupRoutes(app)

	app.Listen(":3000")
}
