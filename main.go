package main

import (
	"github.com/gofiber/fiber/v2"

	"ambassador/src/database"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, hot reload!")
	})

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
