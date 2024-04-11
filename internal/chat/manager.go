package chat

import (
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (s *Socket) NewClient(userId uint, conn *websocket.Conn) *Client {
	return &Client{
		ID:         userId,
		Connection: conn,
	}
}

func (s *Socket) AddConnection(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(client.ID)
	s.clients[client] = true
}

func (s *Socket) RemoveConnection(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client)

}

func (s *Socket) ListenForMessages(r *http.Request, client *Client, chatService *ChatService) {
	defer func() {
		s.RemoveConnection(client)
	}()

	for {
		var msg Message
		msg.Sender = client.ID
		err := wsjson.Read(r.Context(), client.Connection, &msg)
		if err != nil {
			log.Printf("Error reading json: %v", err)
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("Normal closure")
			} else {
				log.Println("Unexpected closure or error")
			}
			break
		}
		if msg.Content != "" {
			msg, err = chatService.StoreMessage(msg)
			if err != nil {
				log.Printf("Error storing message: %v", err)
				break
			}
		}
		s.broadcast <- msg

	}

	client.Connection.Close(websocket.StatusNormalClosure, "Closing connection")
}
