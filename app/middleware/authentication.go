package middleware

import (
	"github.com/gofiber/fiber/v2"
	"social-app/app/dto"
	"social-app/pkg/security"
	"strings"
)

func AuthenticationMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDto{Type: "unauthorized_error", Errors: "requires a token to access this"})
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := security.VerifyToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponseDto{Type: "unauthorized_error", Errors: err.Error()})
	}

	ctx.Locals("userId", claims["id"])
	return ctx.Next()
}
