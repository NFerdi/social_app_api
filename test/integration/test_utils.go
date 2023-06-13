package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	dto2 "social-app/internal/dto"
	"social-app/pkg/database"
)

func getDatabaseConnection() *gorm.DB {
	databaseConnection, err := database.InitMysql()
	if err != nil {
		logrus.Fatalf("Error to open mysql database connection: %v", err)
	}

	return databaseConnection
}

func getBearerToken(app *fiber.App) (string, error) {
	signupDto := dto2.SignupDto{
		Username: "test",
		FullName: "testing",
		Email:    "test@test.com",
		Password: "test",
	}
	loginDto := dto2.LoginDto{
		Email:    "test@test.com",
		Password: "test",
	}
	requestBodySignup, err := json.Marshal(signupDto)
	if err != nil {
		return "", err
	}
	requestBodyLogin, err := json.Marshal(loginDto)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(requestBodySignup))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = app.Test(req)
	if err != nil {
		return "", err
	}

	reqLogin, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBodyLogin))
	if err != nil {
		return "", err
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	resLogin, err := app.Test(reqLogin)
	if err != nil {
		return "", err
	}

	response, err := GetResponseSuccess(resLogin)
	if err != nil {
		return "", err
	}

	token, ok := response.Data.(string)
	if !ok {
		return "", errors.New("token not string")
	}

	return token, nil
}

func getBodyRequestMultipart(fileName string) (*bytes.Buffer, *multipart.Writer, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	rootDir := filepath.Dir(filepath.Dir(currentDir))
	filePath := path.Join(rootDir, "uploads", "default", fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return nil, nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		writer.Close()
		return nil, nil, err
	}

	if fileName == "default_post.jpg" {
		if err = writer.WriteField("caption", "test"); err != nil {
			writer.Close()
			return nil, nil, err
		}
	} else if fileName == "default_avatar.png" {
		if err = writer.WriteField("bio", "test"); err != nil {
			writer.Close()
			return nil, nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, nil, err
	}

	return body, writer, nil
}

func GetResponseSuccess(res *http.Response) (dto2.SuccessResponseDto, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return dto2.SuccessResponseDto{}, err
	}

	var responseDto dto2.SuccessResponseDto
	if err = json.Unmarshal(body, &responseDto); err != nil {
		return dto2.SuccessResponseDto{}, err
	}

	return responseDto, nil
}

func GetResponseError(res *http.Response) (dto2.ErrorResponseDto, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return dto2.ErrorResponseDto{}, err
	}
	res.Body.Close()

	var responseDto dto2.ErrorResponseDto
	if err = json.Unmarshal(body, &responseDto); err != nil {
		return dto2.ErrorResponseDto{}, err
	}

	return responseDto, nil
}

func DeleteAllData(db *gorm.DB) {
	if err := db.Exec("DELETE FROM user_liked_posts").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM user_viewed_posts").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM user_followers").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM posts").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		panic(err)
	}
}
