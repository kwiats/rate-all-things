package auth

import "mime/multipart"

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInDTO struct {
	FirstName      string                `json:"firstname"`
	LastName       string                `json:"lastname"`
	ProfilePicture *multipart.FileHeader `json:"-"`
	Username       string                `json:"username"`
	Email          string                `json:"email"`
	Password       string                `json:"password"`
}
type UserDTO struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	Email          string `json:"email,omitempty"`
}

type LoginDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type Token struct {
	Token string `json:"token"`
}
