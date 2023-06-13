package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"social-app/internal/dto"
)

func CsrfMiddleware(ctx *fiber.Ctx) error {
	if ctx.Method() == http.MethodPost ||
		ctx.Method() == http.MethodPatch ||
		ctx.Method() == http.MethodPut ||
		ctx.Method() == http.MethodDelete {

		token := ctx.Get("X-CSRF-Token")
		if token != "" {
			token = ctx.FormValue("_csrf")
		}

		sessionToken, ok := ctx.Locals("csrfToken").(string)
		if !ok {
			return ctx.Status(fiber.StatusForbidden).JSON(dto.ErrorResponseDto{Type: "", Errors: "invalid csrf token"})
		}
		if token != sessionToken {
			return ctx.Status(fiber.StatusForbidden).JSON(dto.ErrorResponseDto{Type: "", Errors: "invalid csrf token"})
		}
	}

	return ctx.Next()
}
