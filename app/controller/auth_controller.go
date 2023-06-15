package controller

import (
	"github.com/gofiber/fiber/v2"
	"social-app/app/dto"
	"social-app/app/service"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
)

type AuthController interface {
	Signup(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type AuthControllerStruct struct {
	Service service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerStruct{authService}
}

func (c *AuthControllerStruct) Signup(ctx *fiber.Ctx) error {
	signupDto := new(dto.SignupDto)

	if err := ctx.BodyParser(signupDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "validation_error",
			Errors: "invalid request data",
		})
	}

	if err := util.Validate(signupDto); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "validation_error",
			Errors: err,
		})
	}

	if err := c.Service.Signup(signupDto); err != nil {
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
		Message: "successfully registered account",
		Data:    nil,
	})
}

func (c *AuthControllerStruct) Login(ctx *fiber.Ctx) error {
	loginDto := new(dto.LoginDto)

	if err := ctx.BodyParser(loginDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "validation_error",
			Errors: "invalid request data",
		})
	}

	if err := util.Validate(loginDto); len(err) > 0 {

		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
			Type:   "validation_error",
			Errors: err,
		})
	}

	token, err := c.Service.Login(loginDto)
	if err != nil {
		if serviceErr, ok := err.(*customError.ServiceError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   serviceErr.Type,
				Errors: serviceErr.Errors,
			})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponseDto{
				Type:   "service_error",
				Errors: err.Error(),
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.SuccessResponseDto{Message: "successfully logged in", Data: token})
}
