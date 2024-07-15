package routes

import (
	"blogappgolang/controllers"
	"blogappgolang/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRoutes(group fiber.Router) {
	userRoute := group.Group("/user")

	userRoute.Post("/", controllers.CreateUser)
	userRoute.Post("/login", controllers.LoginUser)
	userRoute.Post("/logout", middleware.JWTProtected(), controllers.LogoutUser)
	userRoute.Get("/user-info", middleware.JWTProtected(), controllers.UserInfo)
}
