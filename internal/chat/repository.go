package chat

import (
	"time"
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
func (r *ChatRepository) GetUserChats(userID uint) ([]*model.Chat, error) {
	var chats []*model.Chat

	err := r.db.Preload("Users").Joins("JOIN chat_users on chat_users.chat_id = chats.id").Where("chat_users.user_id = ?", userID).Find(&chats).Error
	if err != nil {
		return nil, err
	}

	return chats, nil
}
func (r *ChatRepository) GetLastMessage(chatID uint) (*model.Message, error) {
	var message model.Message

	err := r.db.Where("chat_id = ?", chatID).Order("created_at DESC").First(&message).Error
	if err != nil {
		return &model.Message{}, err
	}

	return &message, nil
}

func (repo *ChatRepository) GetChatByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	result := repo.db.Preload("Users").Preload("Messages").First(&chat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat, nil
}
func (repo *ChatRepository) SignAsReadBy(messageID uint, userID uint) (bool, error) {
	var exists int64
	repo.db.Model(&model.MessageRead{}).Where("message_id = ? AND user_id = ?", messageID, userID).Count(&exists)
	if exists > 0 {
		return true, nil
	}

	messageRead := model.MessageRead{
		UserID:    userID,
		MessageID: messageID,
		ReadAt:    time.Now(),
	}

	if err := repo.db.Create(&messageRead).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (repo *ChatRepository) GetMessagesFromChat(chatID uint, filters QueryParamsFilters) ([]model.Message, error) {
	query := repo.db.Preload("MessageReads").Where("chat_id = ?", chatID)

	if filters.UserID != 0 {
		query = query.Joins("JOIN message_reads ON message_reads.message_id = messages.id").
			Where("message_reads.user_id = ?", filters.UserID)
	}

	if filters.SenderID != 0 {
		query = query.Where("sender_id = ?", filters.SenderID)
	}

	if filters.StartDate != "" {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != "" {
		query = query.Where("created_at <= ?", filters.EndDate)
	}

	if filters.ContentLike != "" {
		query = query.Where("content LIKE ?", "%"+filters.ContentLike+"%")
	}

	if filters.PageSize > 0 {
		offset := (filters.Page - 1) * filters.PageSize
		query = query.Offset(offset).Limit(filters.PageSize)
	}

	if filters.OrderBy != "" {
		query = query.Order(filters.OrderBy)
	} else {
		query = query.Order("created_at desc")
	}

	var messages []model.Message
	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo *ChatRepository) StoreMessage(message *model.Message) (*model.Message, error) {
	if err := repo.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
