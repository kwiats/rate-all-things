package auth

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kwiats/rate-all-things/pkg/schema"
)

const signingKey = "tajny_kod_do_hashowania_hasła"

func NewAccessToken(user schema.UserDTO) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(60 * time.Minute).Unix()
	claims["name"] = user.Username
	claims["sub"] = user.Id

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckValidationToken(fullToken string) (bool, error) {
	token := strings.TrimPrefix(fullToken, "Bearer ")
	if token == "" {
		log.Printf("Token needed to be provided.")
		return false, errors.New("token not provided")
	}
	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	}

	parsedToken, err := jwt.Parse(token, keyFunc)

	if err != nil {
		log.Printf("Error parsing token: %v\n", err)
		return false, err
	}

	if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return true, nil
	}

	return false, errors.New("invalid token")
}
