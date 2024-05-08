package chat

import (
	"tit/internal/db/model"
	"tit/internal/user"
)

func convertToChatDTO(chat *model.Chat) *ChatDTO {
	userDTOs := make([]user.UserDTO, len(chat.Users))
	for i, userModel := range chat.Users {
		userDTOs[i] = user.UserDTO{Id: userModel.ID, Username: userModel.Username}
	}
	return &ChatDTO{ID: chat.ID, Users: userDTOs}
}
