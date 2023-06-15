package integration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"social-app/app/route"
	"social-app/pkg/config"
	"social-app/pkg/logger"
	"testing"
)

var (
	app *fiber.App
	db  *gorm.DB
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {
	logsPath := filepath.Join("..", "..", "logs")
	if err := logger.InitLogger(logsPath); err != nil {
		logrus.Fatalf("Error creating log directory: %v", err)
	}

	if err := config.InitConfig(filepath.Join("..", "..")); err != nil {
		logrus.Fatalf("Error load file env file: %v", err)
	}

	db = getDatabaseConnection()
	app = setupAppTest()
	DeleteAllData(db)

	logrus.Info("Setup Complete")
}

func teardown() {
	DeleteAllData(db)

	logrus.Info("Teardown complete")
}

func setupAppTest() *fiber.App {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")
	route.MainRoute(apiV1, db)

	return app
}
