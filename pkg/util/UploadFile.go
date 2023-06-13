package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadImage(ctx *fiber.Ctx, folderName string, fileName string) (*multipart.FileHeader, string, error) {
	file, err := ctx.FormFile(fileName)
	if err != nil {
		return nil, "", err
	}

	if err := os.MkdirAll(fmt.Sprintf("./uploads/%s", folderName), os.ModePerm); err != nil {
		return nil, "", err
	}

	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	newFileName := fmt.Sprintf("%d_%s", rand.Intn(9000)+1000, file.Filename)
	dst, err := os.Create(filepath.Join(fmt.Sprintf("./uploads/%s", folderName), newFileName))
	if err != nil {
		return nil, "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, "", err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(wd, "uploads", folderName, newFileName)

	return file, filePath, nil
}

func Upload(file *multipart.FileHeader, folderName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	rootDir := os.Args[1]
	err = os.MkdirAll(filepath.Join(rootDir, "uploads", folderName), 0755)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s", rand.Intn(9000)+1000, file.Filename)

	filePath := filepath.Join(rootDir, "uploads", folderName, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}
