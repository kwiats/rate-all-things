package socket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"tit/internal/db/model"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Socket struct {
	clients   map[*model.Client]bool
	broadcast chan model.Message
	mu        sync.RWMutex
}

func NewSocket() *Socket {
	return &Socket{
		clients:   make(map[*model.Client]bool),
		broadcast: make(chan model.Message),
	}
}


func (s *Socket) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error accepting connection: ", err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	client := s.NewClient(conn)

	s.AddConnection(client)

	for {
		var msg model.Message
		err := wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			log.Println("Error reading json: ", err)
			break
		}
		s.broadcast <- msg
	}

	s.RemoveConnection(client)

	conn.Close(websocket.StatusNormalClosure, "Close connection")
}

func (s *Socket) HandleMessages() {
	for {
		msg := <-s.broadcast
		log.Print(msg)
		s.mu.RLock()
		for client := range s.clients {
			switch msg.Reciver.Type {
			case "user":
				if msg.Reciver.Id == client.Id {
					s.sendMessage(client, &msg)
				}
			case "group":
				if msg.Reciver.Id == client.Id {
					s.sendMessage(client, &msg)
				}
			default:
				s.sendMessage(client, &msg)
			}
		}
		s.mu.RUnlock()
	}
}

func (s *Socket) sendMessage(client *model.Client, msg *model.Message) {
	err := wsjson.Write(context.Background(), client.Connection, msg)
	if err != nil {
		log.Printf("error: %v", err)
		client.Connection.Close(websocket.StatusInternalError, "error writing message")
		s.RemoveConnection(client)
	}
}
