package service

import (
	"errors"
	"fmt"

	"github.com/kwiats/rate-all-things/internal/auth/repository"
	"github.com/kwiats/rate-all-things/pkg/auth"
	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
	utils "github.com/kwiats/rate-all-things/pkg/user"
)

type AuthService struct {
	repository repository.IAuthRepository
}
type IAuthService interface {
	CreateUser() (*model.User, error)
}

func NewAuthService(repository repository.IAuthRepository) *AuthService {
	return &AuthService{repository: repository}
}

func (service *AuthService) LoginUser(user schema.LoginUserDTO) (string, error) {

	userFromDB, err := service.repository.GetUserByUsername(user.Username)
	if err != nil {
		return "", errors.New("cannot found user")
	}

	isLogged := utils.VerifyUserPassword(user.Password, userFromDB.Password)
	if !isLogged {
		return "", errors.New("invalid credentials")
	}

	userDTO := schema.UserDTO{
		Id:       userFromDB.ID,
		Username: userFromDB.Username,
	}
	accessToken, err := auth.NewAccessToken(userDTO)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}
