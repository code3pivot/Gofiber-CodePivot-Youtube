package main

import (
	config "blogappgolang/config"

	"blogappgolang/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDb()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello FIber")
	})

	api := app.Group("/api/v1")
	routes.SetUpUserRoutes(api)
	routes.SetUpCategoryRoutes(api)
	routes.SetUpBlogRoutes(api)

	app.Listen(":8000")

}
