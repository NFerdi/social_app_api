package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"social-app/app/middleware"
	"social-app/app/route"
	"social-app/pkg/config"
	"social-app/pkg/database"
	"social-app/pkg/logger"
	"sync"
	"syscall"
)

func main() {
	if err := logger.InitLogger("logs"); err != nil {
		logrus.Fatalf("Error creating log directory: %v", err)
	}

	if err := config.InitConfig("./"); err != nil {
		logrus.Fatalf("Error load file env file: %v", err)
	}

	app := fiber.New()
	appPort := os.Getenv("APP_PORT")
	db, err := database.InitMysql()

	if err != nil {
		logrus.Fatalf("Error to open mysql database connection: %v", err)
	}

	app.Use(
		middleware.LoggerMiddleware,
	)

	apiV1 := app.Group("/api/v1")

	route.MainRoute(apiV1, db)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := app.Listen(fmt.Sprintf(":%s", appPort)); err != nil {
			logrus.Fatalf("Error while running server : %v", err)
		}
	}()

	logrus.Infof("Server successfully running on port %s", appPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server")

	if err := app.Shutdown(); err != nil {
		logrus.Errorf("Error shutting down server: %s", err.Error())
	}

	wg.Wait()

	logrus.Info("Server stopped")
}
