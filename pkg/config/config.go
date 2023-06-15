package config

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"path"
)

func InitConfig(envPath string) error {
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

	if err := godotenv.Load(path.Join(envPath, envFileName)); err != nil {
		return err
	}

	return nil
}
