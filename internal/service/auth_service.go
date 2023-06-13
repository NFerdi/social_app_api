package service

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"social-app/internal/dto"
	"social-app/internal/entity"
	customError "social-app/pkg/error"
	"social-app/pkg/util"
)

type AuthService interface {
	Signup(request *dto.SignupDto) error
	Login(request *dto.LoginDto) (string, error)
}

type AuthServiceStruct struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &AuthServiceStruct{db: db}
}

func (s *AuthServiceStruct) Signup(request *dto.SignupDto) error {
	var user entity.User

	if err := s.db.Where("email = ? OR username = ?", request.Email, request.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

			if err != nil {
				log.Errorf("Bcrypt error: %s", err.Error())
				return &customError.ServiceError{Type: "encrypted_password_error", Errors: err.Error()}
			}

			request.Password = string(hashedPassword)

			user = entity.User{
				Username: request.Username,
				FullName: request.FullName,
				Email:    request.Email,
				Password: request.Password,
			}

			if err := s.db.Create(&user).Error; err != nil {
				log.Errorf("Database error: %s", err.Error())
				return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
			}

			return nil
		}

		log.Errorf("Database error: %s", err.Error())
		return &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	return &customError.ServiceError{Type: "value_exist", Errors: "email or username already exist"}
}

func (s *AuthServiceStruct) Login(request *dto.LoginDto) (string, error) {
	var user entity.User

	if err := s.db.Where("email = ? OR username = ?", request.Email, request.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", &customError.ServiceError{Type: "value_doesnt_exist", Errors: "account could not be found with this email"}
		}

		log.Errorf("Database error: %s", err.Error())
		return "", &customError.ServiceError{Type: "database_error", Errors: err.Error()}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", &customError.ServiceError{Type: "authentication_error", Errors: "password don't match"}
		}

		log.Errorf("Bcrypt error: %s", err.Error())
		return "", &customError.ServiceError{Type: "bcrypt_error", Errors: err.Error()}
	}

	payloadToken := dto.JwtResponse{
		Id:       user.ID,
		Username: user.Username,
	}
	token, err := util.GenerateToken(payloadToken)
	if err != nil {
		log.Errorf("JWT error: %s", err.Error())
		return "", &customError.ServiceError{Type: "jsonwebtoken_error", Errors: err.Error()}
	}

	return token, nil
}
