package util

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)

func InitConfig() error {
	var envFileName string
	env := flag.String("env", "development", "Environment mode: development, test, production")
	flag.Parse()

	switch *env {
	case "test":
		envFileName = "test.env"
	case "development":
		envFileName = "dev.env"
	default:
		logrus.Fatalf("Invalid environment mode")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("failed to get root directory: %s", err.Error())
		return err
	}
	rootDir := filepath.Dir(filepath.Dir(currentDir))

	if err := godotenv.Load(path.Join(rootDir, envFileName)); err != nil {
		return err
	}

	return nil
}
