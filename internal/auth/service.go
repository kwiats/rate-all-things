package auth

import (
	"errors"
	"fmt"
	"tit/internal/db/model"
	"tit/pkg/util"
)

type AuthService struct {
	repository IAuthRepository
}

type IAuthRepository interface {
	GetUserByUsername(string) (*model.User, error)
	SignIn(*model.User) (*model.User, error)
}

func NewAuthService(repository IAuthRepository) *AuthService {
	return &AuthService{repository: repository}
}

func (service *AuthService) LoginUser(user LoginUserDTO) (string, error) {

	userFromDB, err := service.repository.GetUserByUsername(user.Username)
	if err != nil {
		return "", errors.New("cannot found user")
	}

	isLogged := util.VerifyUserPassword(user.Password, []byte(userFromDB.Password))
	if !isLogged {
		return "", errors.New("invalid credentials")
	}

	userDTO := UserDTO{
		Id:       userFromDB.ID,
		Username: userFromDB.Username,
	}
	accessToken, err := NewAccessToken(userDTO)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (service *AuthService) RegisterUser(userDTO SignInDTO) (*UserDTO, error) {
	hashedPassword, err := util.GeneratePasswordFromString(userDTO.Password)
	if err != nil {
		return &UserDTO{}, err
	}
	
	userDAO := model.User{
		Username: userDTO.Username,
		Email:    userDTO.Email,
		Password: string(*hashedPassword),
	}

	user, err := service.repository.SignIn(&userDAO)
	if err != nil {
		return &UserDTO{}, err
	}
	return &UserDTO{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
