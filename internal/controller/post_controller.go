package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"social-app/internal/dto"
	"social-app/internal/service"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
	"strconv"
)

type PostController interface {
	CreatePost(ctx *fiber.Ctx) error
	WatchPost(ctx *fiber.Ctx) error
	GetUserFeed(ctx *fiber.Ctx) error
	GetViewersOnPost(ctx *fiber.Ctx) error

	LikePost(ctx *fiber.Ctx) error
	UnLikePost(ctx *fiber.Ctx) error
	GetUserLikedPost(ctx *fiber.Ctx) error
}

type PostControllerStruct struct {
	service service.PostService
}

func NewPostController(postService service.PostService) PostController {
	return &PostControllerStruct{postService}
}

func (c *PostControllerStruct) CreatePost(ctx *fiber.Ctx) error {
	userID := uint(uint64(ctx.Locals("userId").(float64)))
	createPostDto := new(dto.CreatePostDto)

	if err := ctx.BodyParser(createPostDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	if err := util.Validate(createPostDto); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: err})
	}

	if err := c.service.CreatePost(ctx, userID, *createPostDto); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: "successfully added post", Data: nil})
}

func (c *PostControllerStruct) WatchPost(ctx *fiber.Ctx) error {
	userID := uint(uint64(ctx.Locals("userId").(float64)))
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err.Error()})
	}

	if err := c.service.ViewPost(userID, uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: fmt.Sprintf("successfully to see post with id %s", stringPostID), Data: nil})
}

func (c *PostControllerStruct) GetUserFeed(ctx *fiber.Ctx) error {
	userID := uint(uint64(ctx.Locals("userId").(float64)))

	posts, err := c.service.GetUserFeed(userID)

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: "successfully to get feed post", Data: posts})
}

func (c *PostControllerStruct) GetViewersOnPost(ctx *fiber.Ctx) error {
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)

	users, err := c.service.GetViewersOnPost(uint(postID))

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: fmt.Sprintf("successfully to get viewers from post id %s", stringPostID), Data: users})
}

func (c *PostControllerStruct) LikePost(ctx *fiber.Ctx) error {
	userID := uint(uint64(ctx.Locals("userId").(float64)))
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err.Error()})
	}

	if err := c.service.LikePost(userID, uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: fmt.Sprintf("successfully to like a post with id %s", stringPostID), Data: nil})
}

func (c *PostControllerStruct) UnLikePost(ctx *fiber.Ctx) error {
	userID := uint(uint64(ctx.Locals("userId").(float64)))
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err.Error()})
	}

	if err := c.service.UnLikePost(userID, uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{Message: fmt.Sprintf("successfully to unlike a post with id %s", stringPostID), Data: nil})
}

func (c *PostControllerStruct) GetUserLikedPost(ctx *fiber.Ctx) error {
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)

	users, err := c.service.GetUserWhoLikedPost(uint(postID))

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: fmt.Sprintf("managed to fetch user data who liked posts with id %s", stringPostID), Data: users})
}
