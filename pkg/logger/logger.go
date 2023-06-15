package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitLogger(basePath string) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	infoFile, err := os.OpenFile(filepath.Join(basePath, "/info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	debugFile, err := os.OpenFile(filepath.Join(basePath, "/debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	warnFile, err := os.OpenFile(filepath.Join(basePath, "/warn.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorFile, err := os.OpenFile(filepath.Join(basePath, "/error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fatalFile, err := os.OpenFile(filepath.Join(basePath, "/fatal.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
