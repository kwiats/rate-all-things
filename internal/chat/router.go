package chat

import (
	"sync"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func HandleWebSocketRouter(db *gorm.DB, router *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	chatRepo := NewChatRepository(db)
	chatService := NewChatService(chatRepo)
	websocket := NewSocket()

	router.HandleFunc("/ws/", HandleConnection(chatService, websocket))

}

func HandleChatRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()
	router := r.PathPrefix("/api/chat").Subrouter()

	chatRepo := NewChatRepository(db)
	chatService := NewChatService(chatRepo)

	router.HandleFunc("/{id}/messages", HandleGetMessages(chatService)).Methods("GET")
	router.HandleFunc("/", HandleGetChats(chatService)).Methods("GET")

}
