package schema

import "github.com/golang-jwt/jwt/v5"

type UserCreateDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserDTO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.Claims
}
