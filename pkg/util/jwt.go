package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"social-app/internal/dto"
	"time"
)

func GenerateToken(payload dto.JwtResponse) (string, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = payload.Id
	claims["username"] = payload.Username
	claims["exp"] = time.Now().AddDate(0, 1, 0).Unix()

	signedToken, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid method signature: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid tokens")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claim")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	if time.Now().After(expirationTime) {
		return nil, errors.New("tokens have expired")
	}

	return claims, nil
}
