package chat

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"tit/internal/auth"
	"tit/pkg/common"

	"github.com/gorilla/mux"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func HandleConnection(chatService *ChatService, socket *Socket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			return
		}

		defer func() {
			if err := conn.Close(websocket.StatusInternalError, "An error occurred"); err != nil {
				log.Print(err)
			}
		}()

		var initialMessage map[string]string
		if err := wsjson.Read(r.Context(), conn, &initialMessage); err != nil {
			log.Printf("Error reading initial message: %v", err)
			return
		}

		token, ok := initialMessage["token"]
		if !ok || token == "" {
			log.Printf("No token found in initial message or token is empty")
			if err := conn.Close(websocket.StatusBadGateway, "Unauthorized - Token not provided or empty"); err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			log.Printf("Token is malformed: token contains an invalid number of segments")
			if err := conn.Close(websocket.StatusBadGateway, "Unauthorized - Malformed token"); err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			return
		}

		userId, err := auth.GetUserIDFromToken(token)
		if err != nil {
			log.Printf("Error extracting user ID from token: %v", err)
			if err := conn.Close(websocket.StatusBadGateway, "Unauthorized - Invalid token"); err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			return
		}

		go socket.HandleMessages(chatService)

		client := socket.NewClient(userId, conn)

		socket.AddConnection(client)
		socket.ListenForMessages(r, client, chatService)
	}
}

func HandleGetMessages(chatService *ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		chatID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		filters := QueryParamsFilters{
			Page:        common.ParseQueryParamInt(r, "page", 1),
			PageSize:    common.ParseQueryParamInt(r, "page_size", 10),
			OrderBy:     r.URL.Query().Get("order_by"),
			StartDate:   r.URL.Query().Get("start_date"),
			EndDate:     r.URL.Query().Get("end_date"),
			ContentLike: r.URL.Query().Get("content_like"),
			SenderID:    common.ParseQueryParamUint(r, "sender_id", 0),
		}
		messages, err := chatService.GetMessagesFromChat(uint(chatID), filters)
		if err != nil {
			return
		}

		common.WriteJSON(w, http.StatusOK, messages)

	}
}

func HandleGetChats(chatService *ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := r.URL.Query().Get("userId")
		if userIdStr == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		userId, err := strconv.ParseUint(userIdStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid User ID", http.StatusBadRequest)
			return
		}

		chats, err := chatService.GetUserChats(uint(userId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		common.WriteJSON(w, http.StatusOK, chats)
	}
}
