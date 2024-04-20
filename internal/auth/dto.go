package auth

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserDTO struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	Email          string `json:"email,omitempty"`
}

type Token struct {
	Token string `json:"token"`
}
