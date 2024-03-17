package userUtils

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordFromString(password string) (*[]byte, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return &hashedPassword, nil
}

func VerifyUserPassword(password string, hashedPassword []byte) bool {
	passwordBytes := []byte(password)
	if err := bcrypt.CompareHashAndPassword(hashedPassword, passwordBytes); err != nil {
		return false
	}
	return true
}
