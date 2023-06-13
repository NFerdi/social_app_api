package util

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitLogger() error {
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("failed to get root directory: %s", err.Error())
		return err
	}
	rootDir := filepath.Dir(filepath.Dir(currentDir))

	logsPath := filepath.Join(rootDir, "logs")

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	infoFile, err := os.OpenFile(filepath.Join(logsPath, "info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	debugFile, err := os.OpenFile(filepath.Join(logsPath, "debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	warnFile, err := os.OpenFile(filepath.Join(logsPath, "warn.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorFile, err := os.OpenFile(filepath.Join(logsPath, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fatalFile, err := os.OpenFile(filepath.Join(logsPath, "fatal.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  infoFile,
			logrus.DebugLevel: debugFile,
			logrus.WarnLevel:  warnFile,
			logrus.ErrorLevel: errorFile,
			logrus.FatalLevel: fatalFile,
		},
		&logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	))

	return nil
}
