package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/welcome", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Nice!")
    })

    app.Listen(":3000")
}