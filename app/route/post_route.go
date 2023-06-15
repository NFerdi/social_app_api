package route

import (
	"github.com/gofiber/fiber/v2"
	"social-app/app/controller"
	"social-app/app/middleware"
)

func PostRoute(router fiber.Router, controller controller.PostController) {
	router.Post("/", middleware.AuthenticationMiddleware, controller.CreatePost)
	router.Get("/", middleware.AuthenticationMiddleware, controller.GetUserFeed)

	router.Post("/:postID/view", middleware.AuthenticationMiddleware, controller.WatchPost)
	router.Get("/:postID/viewers", middleware.AuthenticationMiddleware, controller.GetViewersOnPost)

	router.Post("/:postID/like", middleware.AuthenticationMiddleware, controller.LikePost)
	router.Post("/:postID/unlike", middleware.AuthenticationMiddleware, controller.UnLikePost)
	router.Get("/:postID/likes", middleware.AuthenticationMiddleware, controller.GetUserLikedPost)

	router.Post("/:postID/comment", middleware.AuthenticationMiddleware, controller.CreateCommentPost)
	router.Get("/:postID/comment", middleware.AuthenticationMiddleware, controller.GetPostComments)
	router.Delete("/:postID/comment/:commentID", middleware.AuthenticationMiddleware, controller.DeletePostComment)
}
