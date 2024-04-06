package main

import (
	"github.com/gofiber/fiber/v2"

	"ambassador/src/database"
	"ambassador/src/routes"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()

	routes.Setup(app)

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
