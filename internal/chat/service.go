package chat

import (
	"errors"
	"log"
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
	FindChatForUserIds(userIds []uint) (*model.Chat, error)
	GetUserChats(uint) ([]*model.Chat, error)
	SignAsReadBy(uint, uint) (bool, error)
	GetLastMessage(uint) (*model.Message, error)
	GetMessagesFromChat(uint, QueryParamsFilters) ([]model.Message, error)
}

func NewChatService(repository IChatRepository) *ChatService {
	return &ChatService{repository: repository}
}

func (service *ChatService) GetOrCreateChat(chat ChatDTO) (*ChatDTO, error) {
	users := make([]model.User, len(chat.Users))
	userIds := make([]uint, len(chat.Users))
	
	for i, userDTO := range chat.Users {
		users[i] = model.User{Model: &gorm.Model{ID: userDTO.Id}, Username: userDTO.Username}
		userIds[i] = userDTO.Id
	}

	existingChat, err := service.repository.FindChatForUserIds(userIds)
	if err == nil && existingChat != nil {
		return convertToChatDTO(existingChat), nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	chatdao := &model.Chat{Users: users}
	chatdao, err = service.repository.CreateChat(chatdao)
	if err != nil {
		return nil, err
	}

	return convertToChatDTO(chatdao), nil
}

func (service *ChatService) GetUserChats(userId uint) ([]ChatsList, error) {
	chatsdao, err := service.repository.GetUserChats(userId)
	if err != nil {
		return []ChatsList{}, err
	}
	chatsDTO := make([]ChatsList, len(chatsdao))
	for i, chatdao := range chatsdao {
		usersDTO := make([]user.UserDTO, len(chatdao.Users))

		for i, userdao := range chatdao.Users {
			usersDTO[i] = user.UserDTO{
				Id:       userdao.ID,
				Username: userdao.Username,
			}
		}
		message, _ := service.GetLastMessage(chatdao.ID)
		chatsDTO[i] = ChatsList{
			ID:          chatdao.ID,
			Users:       usersDTO,
			LastMessage: message,
		}

	}

	return chatsDTO, nil
}
func (service *ChatService) GetLastMessage(chatID uint) (Message, error) {
	messagedao, err := service.repository.GetLastMessage(chatID)
	if err != nil {
		return Message{}, err
	}
	return Message{
		ID:        messagedao.ID,
		CreatedAt: messagedao.CreatedAt,
		Content:   messagedao.Content,
		Sender:    messagedao.SenderID,
	}, nil
}

func (service *ChatService) GetChatByID(id uint) (*ChatDTO, error) {
	chatdao, err := service.repository.GetChatByID(id)
	if err != nil {
		return &ChatDTO{}, err
	}

	userDTOs := make([]user.UserDTO, len(chatdao.Users))
	for i, userModel := range chatdao.Users {
		userDTOs[i] = user.UserDTO{Id: userModel.ID}
	}
	return &ChatDTO{Users: userDTOs}, nil
}

func (service *ChatService) GetMessagesFromChat(id uint, filters QueryParamsFilters) ([]MessageDTO, error) {
	messagesDAO, err := service.repository.GetMessagesFromChat(id, filters)
	if err != nil {
		return nil, err
	}

	var messages []MessageDTO
	for _, dao := range messagesDAO {
		readByUsers := []MessageRead{}
		for _, read := range dao.MessageReads {
			readByUsers = append(readByUsers, MessageRead{
				UserID: read.UserID,
				ReadAt: read.ReadAt,
			})
		}

		messages = append(messages, MessageDTO{
			ID:          dao.ID,
			Sender:      dao.SenderID,
			Content:     dao.Content,
			CreatedAt:   dao.CreatedAt,
			Chat:        dao.ChatID,
			ReadByUsers: readByUsers,
		})
	}

	return messages, nil
}

func (service *ChatService) StoreMessage(message Message) (Message, error) {
	mess, err := service.repository.StoreMessage(&model.Message{
		ChatID:   message.Chat,
		SenderID: message.Sender,
		Content:  message.Content,
	})
	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:        mess.ID,
		Chat:      mess.ChatID,
		Sender:    mess.SenderID,
		Content:   mess.Content,
		CreatedAt: mess.CreatedAt,
	}, nil
}

func (service *ChatService) SignAsReadBy(messageID uint, userID uint) bool {
	signed, err := service.repository.SignAsReadBy(messageID, userID)
	if err != nil {
		log.Printf("Error during sign message as read: %v\n", err)
		return false
	}
	return signed
}
