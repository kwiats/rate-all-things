package auth

import (
	"errors"
	"fmt"
	"tit/internal/db/model"
	"tit/pkg/media"
	"tit/pkg/util"
)

type AuthService struct {
	repository IAuthRepository
	storage    *media.MediaService
}

type IAuthRepository interface {
	GetUserByUsername(string) (*model.User, error)
	SignIn(*model.User) (*model.User, error)
}

func NewAuthService(repository IAuthRepository, mediaStorage *media.MediaService) *AuthService {
	return &AuthService{repository: repository, storage: mediaStorage}
}

func (service *AuthService) LoginUser(user LoginUserDTO) (*LoginDTO, error) {
	userFromDB, err := service.repository.GetUserByUsername(user.Username)
	if err != nil {
		// Wrap the error to maintain the underlying error context
		return nil, fmt.Errorf("user not found: %w", err)
	}

	isLogged := util.VerifyUserPassword(user.Password, []byte(userFromDB.Password))
	if !isLogged {
		return nil, errors.New("invalid credentials")
	}

	userDTO := UserDTO{
		Id:             userFromDB.ID,
		Username:       userFromDB.Username,
		FirstName:      userFromDB.FirstName,
		LastName:       userFromDB.LastName,
		ProfilePicture: userFromDB.ProfilePicture,
	}

	accessToken, err := NewAccessToken(userDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	loginDTO := &LoginDTO{Token: accessToken, User: userDTO}
	return loginDTO, nil
}

func (service *AuthService) RegisterUser(userDTO SignInDTO) (*UserDTO, error) {
	hashedPassword, err := util.GeneratePasswordFromString(userDTO.Password)
	if err != nil {
		return &UserDTO{}, err
	}
	var profilePicturePath string
	if userDTO.ProfilePicture != nil {
		profilePicturePath, err = service.storage.SaveFile(userDTO.ProfilePicture, userDTO.Username)
		if err != nil {
			return nil, fmt.Errorf("error saving profile picture: %w", err)
		}
	}

	userDAO := model.User{
		FirstName:      userDTO.FirstName,
		LastName:       userDTO.LastName,
		ProfilePicture: profilePicturePath,
		Username:       userDTO.Username,
		Email:          userDTO.Email,
		Password:       string(*hashedPassword),
	}

	user, err := service.repository.SignIn(&userDAO)
	if err != nil {
		return &UserDTO{}, err
	}

	return &UserDTO{
		Id:             user.ID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}, nil
}
