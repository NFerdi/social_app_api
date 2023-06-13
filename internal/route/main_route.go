package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	controller2 "social-app/internal/controller"
	service2 "social-app/internal/service"
)

func MainRoute(route fiber.Router, db *gorm.DB) {
	authRoute := route.Group("/auth")
	authService := service2.NewAuthService(db)
	authController := controller2.NewAuthController(authService)
	AuthRoute(authRoute, authController)

	userRoute := route.Group("/user")
	userService := service2.NewUserService(db)
	userController := controller2.NewUserController(userService)
	UserRoute(userRoute, userController)

	postRoute := route.Group("/post")
	postService := service2.NewPostService(db)
	postController := controller2.NewPostController(postService)
	PostRoute(postRoute, postController)
}
