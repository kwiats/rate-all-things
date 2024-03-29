package user

import (
	"tit/internal/db/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) DeleteUser(uint) (bool, error) {
	panic("unimplemented")
}

func (repo *UserRepository) UpdateUser(uint, *model.User) (*model.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUsers() ([]*model.User, error) {
	var users []*model.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
