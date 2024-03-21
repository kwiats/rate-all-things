package chat

import (
	"tit/internal/db/model"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (repo *ChatRepository) CreateChat(chat *model.Chat) (*model.Chat, error) {
	if err := repo.db.Create(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

func (repo *ChatRepository) GetChatByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	result := repo.db.Preload("Users").Preload("Messages").First(&chat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat, nil
}

func (repo *ChatRepository) StoreMessage(message *model.Message) (*model.Message, error) {
	if err := repo.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
