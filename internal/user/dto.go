package user

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserCreateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserDTO struct {
	Id             uint   `json:"id"`
	Username       string `json:"username,omitempty"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	Email          string `json:"email,omitempty"`
}
type UserClaimsDTO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.Claims
}
