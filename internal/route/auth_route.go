package route

import (
	"github.com/gofiber/fiber/v2"
	"social-app/internal/controller"
)

func AuthRoute(router fiber.Router, controller controller.AuthController) {

	router.Post("/signup", controller.Signup)
	router.Post("/login", controller.Login)
}
