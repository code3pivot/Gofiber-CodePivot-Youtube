package routes

import (
	"blogappgolang/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUpCategoryRoutes(group fiber.Router) {
	categoryRoute := group.Group("/category")

	categoryRoute.Post("/", controllers.CreateCategory)
	categoryRoute.Get("/", controllers.AllCategory)
	categoryRoute.Get("/:id", controllers.SingleCategoryDetail)
	categoryRoute.Patch("/:id", controllers.UpdateCategory)
	categoryRoute.Delete("/:id", controllers.DeleteCategory)

}
