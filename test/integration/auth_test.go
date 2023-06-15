package integration

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"social-app/app/dto"
	"testing"
)

func TestSignup(t *testing.T) {
	signupDto := dto.SignupDto{
		Username: "test",
		FullName: "testing",
		Email:    "test@test.com",
		Password: "test",
	}
	requestBody, err := json.Marshal(signupDto)
	assert.Nil(t, err)

	t.Run("SuccessSignup", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(requestBody))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "successfully registered account", response.Message)
	})

	t.Run("ErrorWhenEmailAlreadyExist", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(requestBody))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "value_exist", response.Type)
		assert.Equal(t, "email or username already exist", response.Errors)
	})

	t.Run("ErrorWhenEmptyRequest", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "validation_error", response.Type)
		assert.NotNil(t, response.Errors)
	})
}

func TestLogin(t *testing.T) {
	loginDto := dto.LoginDto{
		Email:    "test@test.com",
		Password: "test",
	}
	requestBodyLogin, err := json.Marshal(loginDto)
	assert.Nil(t, err)

	t.Run("SuccessLoginAndReturnToken", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBodyLogin))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)

		response, err := GetResponseSuccess(res)
		assert.Nil(t, err)

		assert.Equal(t, "successfully logged in", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("ErrorWhenEmailNotExist", func(t *testing.T) {
		loginDto.Email = "test1"

		requestBodyLogin, err = json.Marshal(loginDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBodyLogin))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "value_doesnt_exist", response.Type)
		assert.Equal(t, "account could not be found with this email", response.Errors)
	})

	t.Run("ErrorWhenPasswordWrong", func(t *testing.T) {
		loginDto.Email = "test@test.com"
		loginDto.Password = "test1"

		requestBodyLogin, err = json.Marshal(loginDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBodyLogin))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "authentication_error", response.Type)
		assert.Equal(t, "password don't match", response.Errors)
	})

	t.Run("ErrorWhenEmptyRequest", func(t *testing.T) {
		loginDto = dto.LoginDto{}

		requestBodyLogin, err = json.Marshal(loginDto)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		response, err := GetResponseError(res)
		assert.Nil(t, err)

		assert.Equal(t, "validation_error", response.Type)
		assert.NotNil(t, response.Errors)
	})
}
