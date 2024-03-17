package socket

import (
	"fmt"
	"tit/internal/db/model"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

func (s *Socket) NewClient(conn *websocket.Conn) *model.Client {
	return &model.Client{
		Id:         uuid.New().String(),
		Connection: conn,
	}
}

func (s *Socket) AddConnection(client *model.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(client.Id)
	s.clients[client] = true
}

func (s *Socket) RemoveConnection(client *model.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client)

}
