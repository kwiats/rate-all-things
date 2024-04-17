package user

import (
	"mime/multipart"
	"sync"
	"tit/internal/db/model"
	"tit/pkg/media"
	"tit/pkg/util"
)

type UserService struct {
	repository IUserRepository
	storage    media.IBlobStorage
}
type IUserRepository interface {
	CreateUser(*model.User) (*model.User, error)
	GetUserByUsername(string) (*model.User, error)
	GetUserByID(uint) (*model.User, error)
	GetUsers() ([]*model.User, error)
	DeleteUser(uint) (bool, error)
	UpdateUser(uint, *model.User) (*model.User, error)
}

func NewUserService(repository IUserRepository, blobStorage media.IBlobStorage) *UserService {
	return &UserService{repository: repository, storage: blobStorage}
}

func (service *UserService) CreateUser(user UserCreateDTO, profilePicture *multipart.FileHeader) (*UserDTO, error) {
	hashedPassword, _ := util.GeneratePasswordFromString(user.Password)

	pathToProfilePicture, err := service.storage.SaveFile(profilePicture, "")
	if err != nil {
		return nil, err
	}

	userModel := model.User{
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: pathToProfilePicture,
		Password:       string(*hashedPassword),
	}

	createdUser, err := service.repository.CreateUser(&userModel)
	if err != nil {
		return &UserDTO{}, err
	}
	return &UserDTO{
		Id:             createdUser.ID,
		Username:       createdUser.Username,
		Email:          createdUser.Email,
		ProfilePicture: createdUser.ProfilePicture,
	}, nil

}

func (service *UserService) GetUsers() ([]*UserDTO, error) {
	users, err := service.repository.GetUsers()
	if err != nil {
		return nil, err
	}

	users_dto := make([]*UserDTO, 0, len(users))
	channel := make(chan *UserDTO, len(users))

	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(u *model.User) {
			defer wg.Done()
			channel <- &UserDTO{
				Id:       u.ID,
				Username: u.Username,
			}
		}(user)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for _user := range channel {
		users_dto = append(users_dto, _user)
	}

	return users_dto, nil
}
func (service *UserService) GetUserById(userId uint) (*UserDTO, error) {

	user, err := service.repository.GetUserByID(userId)
	if err != nil {
		return &UserDTO{}, err
	}

	return &UserDTO{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (service *UserService) GetUserByUsername(username string) (*UserDTO, error) {

	user, err := service.repository.GetUserByUsername(username)
	if err != nil {
		return &UserDTO{}, err
	}

	return &UserDTO{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
func (service *UserService) DeleteUser(userId uint, force bool) (bool, error) {
	_, err := service.repository.DeleteUser(userId)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (service *UserService) UpdateUser(id uint, userUpdate UserUpdateDTO) (*UserDTO, error) {
	user := model.User{
		Username: userUpdate.Username,
	}
	updatedUser, err := service.repository.UpdateUser(id, &user)
	if err != nil {
		return &UserDTO{}, err
	}

	return &UserDTO{Username: updatedUser.Username}, nil
}
