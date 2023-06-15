package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"social-app/app/dto"
	"social-app/app/service"
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

	CreateCommentPost(ctx *fiber.Ctx) error
	GetPostComments(ctx *fiber.Ctx) error
	DeletePostComment(ctx *fiber.Ctx) error
}

type PostControllerImpl struct {
	service service.PostService
}

func NewPostController(postService service.PostService) PostController {
	return &PostControllerImpl{postService}
}

func (c *PostControllerImpl) CreatePost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	createPostDto := new(dto.CreatePostDto)

	if err := ctx.BodyParser(createPostDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	if err := util.Validate(createPostDto); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: err})
	}

	if err := c.service.CreatePost(ctx, uint(uint64(userID)), *createPostDto); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{
		Message: "successfully added post",
		Data:    nil,
	})
}

func (c *PostControllerImpl) WatchPost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	if err := c.service.ViewPost(uint(uint64(userID)), uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("successfully to see post with id %s", stringPostID),
		Data:    nil,
	})
}

func (c *PostControllerImpl) GetUserFeed(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	posts, err := c.service.GetUserFeed(uint(uint64(userID)))

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{
		Message: "successfully to get feed post",
		Data:    posts,
	})
}

func (c *PostControllerImpl) GetViewersOnPost(ctx *fiber.Ctx) error {
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	users, err := c.service.GetViewersOnPost(uint(postID))

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("successfully to get viewers from post id %s", stringPostID),
		Data:    users,
	})
}

func (c *PostControllerImpl) LikePost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	if err := c.service.LikePost(uint(uint64(userID)), uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("successfully to like a post with id %s", stringPostID),
		Data:    nil,
	})
}

func (c *PostControllerImpl) UnLikePost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	if err := c.service.UnLikePost(uint(uint64(userID)), uint(postID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(http.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("successfully to unlike a post with id %s", stringPostID),
		Data:    nil,
	})
}

func (c *PostControllerImpl) GetUserLikedPost(ctx *fiber.Ctx) error {
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	users, err := c.service.GetUserWhoLikedPost(uint(postID))

	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: serviceErr.Type, Errors: serviceErr.Errors})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "service_error", Errors: err})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("managed to fetch user data who liked posts with id %s", stringPostID),
		Data:    users,
	})
}

func (c *PostControllerImpl) CreateCommentPost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userId").(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: "failed to retrieve user ID from context",
		})
	}

	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	createCommentDto := new(dto.CreateCommentPostDto)
	if err := ctx.BodyParser(createCommentDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: "invalid request data"})
	}

	if err := util.Validate(createCommentDto); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{Type: "validation_error", Errors: err})
	}

	if err := c.service.CreateCommentPost(uint(postID), uint(uint64(userID)), *createCommentDto); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   serviceErr.Type,
				Errors: serviceErr.Errors,
			})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   "service_error",
				Errors: err,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("successfully added comment to post id %s", stringPostID),
		Data:    nil,
	})
}

func (c *PostControllerImpl) GetPostComments(ctx *fiber.Ctx) error {
	stringPostID := ctx.Params("postID")
	postID, err := strconv.ParseUint(stringPostID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	comments, err := c.service.GetPostComments(uint(postID))
	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   serviceErr.Type,
				Errors: serviceErr.Errors,
			})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   "service_error",
				Errors: err,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("managed to get comments on post id %s", stringPostID),
		Data:    comments,
	})
}

func (c *PostControllerImpl) DeletePostComment(ctx *fiber.Ctx) error {
	stringCommentID := ctx.Params("commentID")
	CommentID, err := strconv.ParseUint(stringCommentID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "service_error",
			Errors: err.Error(),
		})
	}

	if err := c.service.DeletePostComment(uint(CommentID)); err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   serviceErr.Type,
				Errors: serviceErr.Errors,
			})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   "service_error",
				Errors: err,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{
		Message: fmt.Sprintf("managed to delete comment in %s", stringCommentID),
		Data:    nil,
	})
}
