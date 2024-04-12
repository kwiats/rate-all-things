package chat

import (
	"context"
	"log"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Socket struct {
	clients   map[*Client]bool
	broadcast chan Message
	mu        sync.RWMutex
}

func NewSocket() *Socket {
	return &Socket{
		clients:   make(map[*Client]bool),
		broadcast: make(chan Message),
	}
}

func (s *Socket) HandleMessages(chatService *ChatService) {
	for {
		msg := <-s.broadcast

		s.mu.RLock()
		chat, err := chatService.GetChatByID(msg.Chat)
		if err != nil {
			log.Printf("error getting chat by ID: %v", err)
			s.mu.RUnlock()
			continue
		}
		userIDs := make(map[uint]struct{})
		for _, user := range chat.Users {
			userIDs[user.Id] = struct{}{}
		}

		for client := range s.clients {
			if _, exists := userIDs[client.ID]; exists {
				go func(client *Client) {
					if err := s.sendMessage(client, &msg); err != nil {
						log.Printf("Error sending message to client %v: %v", client.ID, err)
					}
					chatService.SignAsReadBy(msg.ID, client.ID)
				}(client)
			}
		}

		s.mu.RUnlock()
	}
}

func (s *Socket) sendMessage(client *Client, msg *Message) error {
	err := wsjson.Write(context.Background(), client.Connection, msg)
	if err != nil {
		client.Connection.Close(websocket.StatusInternalError, "error writing message")
		s.RemoveConnection(client)
		return err
	}
	return nil
}
