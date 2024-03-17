package router

import (
	"sync"
	socket "tit/internal/handler/websocket"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func HandleWebsocket(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	websocket := socket.NewSocket()

	r.HandleFunc("/ws", websocket.HandleConnections)
	go websocket.HandleMessages()
}
