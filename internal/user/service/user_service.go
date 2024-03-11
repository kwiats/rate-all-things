package service

import (
	"sync"

	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
	userUtils "github.com/kwiats/rate-all-things/pkg/user"
)

type UserService struct {
	repository IUserRepository
}
type IUserRepository interface {
	CreateUser(*model.User) (*model.User, error)
	GetUserByUsername(string) (*model.User, error)
	GetUserByID(uint) (*model.User, error)
	GetUsers() ([]*model.User, error)
}

func NewUserService(repository IUserRepository) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) CreateUser(user schema.UserCreateDTO) (*schema.UserDTO, error) {
	hashedPassword, _ := userUtils.GeneratePasswordFromString(user.Password)

	userModel := model.User{
		Username: user.Username,
		Password: *hashedPassword,
	}

	createdUser, err := service.repository.CreateUser(&userModel)
	if err != nil {
		return &schema.UserDTO{}, err
	}
	return &schema.UserDTO{
		Id:       createdUser.ID,
		Username: createdUser.Username,
	}, nil

}

func (service *UserService) GetUsers() ([]*schema.UserDTO, error) {
	users, err := service.repository.GetUsers()
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*schema.UserDTO, 0, len(users))
	channel := make(chan *schema.UserDTO, len(users))

	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(u *model.User) {
			defer wg.Done()
			channel <- &schema.UserDTO{
				Id:       u.ID,
				Username: u.Username,
			}
		}(user)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for dto := range channel {
		usersDTO = append(usersDTO, dto)
	}

	return usersDTO, nil
}
func (service *UserService) GetUserById(userId uint) (*schema.UserDTO, error) {

	user, err := service.repository.GetUserByID(userId)
	if err != nil {
		return &schema.UserDTO{}, err
	}

	return &schema.UserDTO{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (service *UserService) GetUserByUsername(username string) (*schema.UserDTO, error) {

	user, err := service.repository.GetUserByUsername(username)
	if err != nil {
		return &schema.UserDTO{}, err
	}

	return &schema.UserDTO{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
func (service *UserService) DeleteUser(userId uint, force bool) (bool, error) {
	if force {
		return false, nil
	}
	return true, nil
}
func (service *UserService) UpdateUser(user schema.UserUpdateDTO) (*schema.UserDTO, error) {
	return nil, nil
}
