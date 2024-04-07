package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"ambassador/src/database"
	"ambassador/src/routes"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost",
	}))

	routes.Setup(app)

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
