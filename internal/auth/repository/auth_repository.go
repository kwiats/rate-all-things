package repository

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

type IAuthRepository interface {
	GetUserByUsername(string) (*model.User, error)
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
