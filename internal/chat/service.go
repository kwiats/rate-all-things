package chat

import (
	"tit/internal/db/model"
	"tit/internal/user"

	"gorm.io/gorm"
)

type ChatService struct {
	repository IChatRepository
}

type IChatRepository interface {
	CreateChat(*model.Chat) (*model.Chat, error)
	StoreMessage(*model.Message) (*model.Message, error)
	GetChatByID(uint) (*model.Chat, error)
}

func NewChatService(repository IChatRepository) *ChatService {
	return &ChatService{repository: repository}
}

func (service *ChatService) CreateChat(chat ChatDTO) (*ChatDTO, error) {
	_users := make([]model.User, len(chat.Users))
	for i, userDTO := range chat.Users {
		_users[i] = model.User{Model: &gorm.Model{ID: userDTO.Id}, Username: userDTO.Username}
	}

	_chat := &model.Chat{
		Users: _users,
	}

	chatdao, err := service.repository.CreateChat(_chat)
	if err != nil {
		return &ChatDTO{}, err
	}

	_userDTOs := make([]user.UserDTO, len(chatdao.Users))
	for i, userModel := range chatdao.Users {
		_userDTOs[i] = user.UserDTO{Id: userModel.ID, Username: userModel.Username}
	}

	return &ChatDTO{Users: _userDTOs}, nil
}

func (service *ChatService) GetChatByID(id uint) (*ChatDTO, error) {
	chatdao, err := service.repository.GetChatByID(id)
	if err != nil {
		return &ChatDTO{}, err
	}

	_userDTOs := make([]user.UserDTO, len(chatdao.Users))
	for i, userModel := range chatdao.Users {
		_userDTOs[i] = user.UserDTO{Id: userModel.ID}
	}
	return &ChatDTO{Users: _userDTOs}, nil
}

func (service *ChatService) StoreMessage(message Message) (*Message, error) {
	mess, err := service.repository.StoreMessage(&model.Message{
		ChatID:   message.Chat,
		SenderID: message.Sender,
		Content:  message.Content,
	})
	if err != nil {
		return &Message{}, err
	}

	return &Message{
		Chat:    mess.ChatID,
		Sender:  mess.SenderID,
		Content: mess.Content}, nil
}
