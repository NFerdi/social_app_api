package route

import (
	"github.com/gofiber/fiber/v2"
	"social-app/internal/controller"
	"social-app/internal/middleware"
)

func PostRoute(router fiber.Router, controller controller.PostController) {
	router.Get("/", middleware.AuthenticationMiddleware, controller.GetUserFeed)
	router.Post("/", middleware.AuthenticationMiddleware, controller.CreatePost)

	router.Post("/:postID/view", middleware.AuthenticationMiddleware, controller.WatchPost)
	router.Get("/:postID/viewers", middleware.AuthenticationMiddleware, controller.GetViewersOnPost)

	router.Post("/:postID/like", middleware.AuthenticationMiddleware, controller.LikePost)
	router.Post("/:postID/unlike", middleware.AuthenticationMiddleware, controller.UnLikePost)
	router.Get("/:postID/likes", middleware.AuthenticationMiddleware, controller.GetUserLikedPost)
}
