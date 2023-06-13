package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()

	err := ctx.Next()

	latency := time.Since(start)
	status := ctx.Response().StatusCode()
	ip := ctx.IP()
	method := ctx.Method()
	path := ctx.Path()

	logrus.Infof("%s | %d | %s | %s | %s | %s", time.Now().Format("15:04:05"), status, latency.String(), ip, method, path)

	return err
}
