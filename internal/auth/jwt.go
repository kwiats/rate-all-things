package auth

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	Username string
	UserID   uint
}

const signingKey = "tajny_kod_do_hashowania_hasła"

func NewAccessToken(user UserDTO) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * 60 * time.Minute).Unix()
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

func GetUserIDFromToken(fullToken string) (uint, error) {
	token := strings.TrimPrefix(fullToken, "Bearer ")
	if token == "" {
		log.Printf("Token needed to be provided.")
		return 0, errors.New("token not provided")
	}

	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	}

	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		log.Printf("Error parsing token: %v\n", err)
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if userID, ok := claims["sub"].(float64); ok {
			return uint(userID), nil
		}
	}

	return 0, errors.New("invalid token")
}
func GetUserInfoFromToken(fullToken string) (UserInfo, bool, error) {
	userInfo := UserInfo{}

	token := strings.TrimPrefix(fullToken, "Bearer ")
	if token == "" {
		log.Printf("Token needed to be provided.")
		return userInfo, false, errors.New("token not provided")
	}

	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	}

	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		log.Printf("Error parsing token: %v\n", err)
		return userInfo, false, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if username, ok := claims["name"].(string); ok {
			userInfo.Username = username
		}
		if userID, ok := claims["sub"].(float64); ok {
			userInfo.UserID = uint(userID)
		}
		return userInfo, true, nil
	}

	return userInfo, false, errors.New("invalid token")
}
