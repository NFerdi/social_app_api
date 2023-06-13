package route

import (
	"github.com/gofiber/fiber/v2"
	"social-app/internal/controller"
	"social-app/internal/middleware"
)

func UserRoute(router fiber.Router, controller controller.UserController) {
	router.Patch("/", middleware.AuthenticationMiddleware, controller.UpdateUserProfile)
	router.Get("/:username", middleware.AuthenticationMiddleware, controller.GetUserByUsername)

	router.Post("/follow", middleware.AuthenticationMiddleware, controller.FollowUser)
	router.Post("/unfollow", middleware.AuthenticationMiddleware, controller.UnfollowUser)

	router.Get("/:username/followers", middleware.AuthenticationMiddleware, controller.GetUserFollowers)
	router.Get("/:username/following", middleware.AuthenticationMiddleware, controller.GetUserFollowing)

	router.Get("/:username/post/uploaded", middleware.AuthenticationMiddleware, controller.GetUserUploadedPosts)
	router.Get("/:username/post/viewed", middleware.AuthenticationMiddleware, controller.GetPostsWatchedByUser)
	router.Get("/:username/post/liked", middleware.AuthenticationMiddleware, controller.GetUserLikedPosts)

}
