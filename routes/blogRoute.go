package routes

import (
	"blogappgolang/controllers"
	"blogappgolang/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpBlogRoutes(group fiber.Router) {
	blogRoute := group.Group("/blog")

	blogRoute.Post("/", controllers.CreateBlog)
	blogRoute.Get("/", controllers.GetAllBlogs)
	blogRoute.Get("/:id", controllers.GetSingleBlogByID)
	blogRoute.Patch("/:id", middleware.JWTProtected(), controllers.UpdateBlog)
	blogRoute.Delete("/:id", middleware.JWTProtected(), controllers.DeleteBlog)

}
