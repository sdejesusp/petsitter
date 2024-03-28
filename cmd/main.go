package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/petsitter/database"
	"github.com/sdejesusp/petsitter/handlers"
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

	setupRoutesWithoutJWT(app)

	jwtSecret := handlers.ReadEnvVariable(handlers.JWTSECRETENV)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
	}))

	// API routes
	setupRoutes(app)

	app.Listen(":3000")
}
