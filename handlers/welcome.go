package handlers

import "github.com/gofiber/fiber/v2"

func WelcomeMessage(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Welcome to the PetSitter API")
}