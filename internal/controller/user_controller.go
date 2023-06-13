package controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"social-app/internal/dto"
	"social-app/internal/service"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
)

type UserController interface {
	GetUserByUsername(ctx *fiber.Ctx) error
	FollowUser(ctx *fiber.Ctx) error
	UnfollowUser(ctx *fiber.Ctx) error
	UpdateUserProfile(ctx *fiber.Ctx) error
	GetUserFollowers(ctx *fiber.Ctx) error
	GetUserFollowing(ctx *fiber.Ctx) error
	GetPostsWatchedByUser(ctx *fiber.Ctx) error
	GetUserUploadedPosts(ctx *fiber.Ctx) error
	GetUserLikedPosts(ctx *fiber.Ctx) error
}

type UserControllerStruct struct {
	service service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerStruct{userService}
}

func (c *UserControllerStruct) GetUserByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	userData, err := c.service.GetUserProfile(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to fetch user based on username", Data: userData})
}

func (c *UserControllerStruct) FollowUser(ctx *fiber.Ctx) error {
	followUser := new(dto.FollowUserDto)
	userId := uint(uint64(ctx.Locals("userId").(float64)))

	if err := ctx.BodyParser(followUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	followerId := uint(uint64(followUser.FollowerId))
	err := c.service.FollowUser(followerId, userId)
	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to follow user", Data: nil})
}

func (c *UserControllerStruct) UnfollowUser(ctx *fiber.Ctx) error {
	followUser := new(dto.FollowUserDto)
	userId := uint(uint64(ctx.Locals("userId").(float64)))

	if err := ctx.BodyParser(followUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	followerId := uint(uint64(followUser.FollowerId))
	err := c.service.UnfollowUser(followerId, userId)
	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: "successfully unfollowed user", Data: nil})
}

func (c *UserControllerStruct) UpdateUserProfile(ctx *fiber.Ctx) error {
	userId := uint(uint64(ctx.Locals("userId").(float64)))
	updateUserProfileDto := new(dto.UpdateUserProfileDto)

	form, err := ctx.MultipartForm()
	if err != nil {

	}

	files := form.File["avatar"]

	if err := ctx.BodyParser(updateUserProfileDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	if err := util.Validate(updateUserProfileDto); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: err})
	}

	if err := c.service.UpdateUserProfile(files, userId, *updateUserProfileDto); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to change the user profile", Data: nil})
}

func (c *UserControllerStruct) GetUserFollowers(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	userData, err := c.service.GetUserFollowers(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to fetch user followers", Data: userData})
}

func (c *UserControllerStruct) GetUserFollowing(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	userData, err := c.service.GetUserFollowing(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to fetch user following", Data: userData})
}

func (c *UserControllerStruct) GetPostsWatchedByUser(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	posts, err := c.service.GetPostsWatchedByUser(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to get posts that have been watched by users", Data: posts})
}

func (c *UserControllerStruct) GetUserUploadedPosts(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	posts, err := c.service.GetUserUploadedPosts(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to fetch all posts uploaded by users", Data: posts})
}

func (c *UserControllerStruct) GetUserLikedPosts(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	posts, err := c.service.GetUserLikedPosts(username)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "managed to get posts that users like", Data: posts})
}
