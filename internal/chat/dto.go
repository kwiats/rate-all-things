package chat

import (
	"time"
	"tit/internal/user"

	"nhooyr.io/websocket"
)

type Client struct {
	ID         uint            `json:"id"`
	Connection *websocket.Conn `json:"-"`
	LastActive time.Time       `json:"last_active,omitempty"`
}

type ChatDTO struct {
	ID       uint           `json:"id"`
	Users    []user.UserDTO `json:"users"`
	Messages []string       `json:"messages"`
}

type ChatsList struct {
	ID          uint           `json:"id"`
	Users       []user.UserDTO `json:"users"`
	LastMessage Message        `json:"last_message,omitempty"`
}

type Message struct {
	ID          uint          `json:"id,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	Sender      uint          `json:"sender,omitempty"`
	Chat        uint          `json:"chat,omitempty"`
	Content     string        `json:"content"`
	ReadByUsers []MessageRead `json:"read_by_users,omitempty"`
}
type MessageDTO struct {
	ID          uint          `json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	Sender      uint          `json:"sender"`
	Chat        uint          `json:"chat"`
	Content     string        `json:"content"`
	ReadByUsers []MessageRead `json:"read_by_users"`
}

type MessageRead struct {
	UserID uint      `json:"user_id"`
	ReadAt time.Time `json:"read_at"`
}

type QueryParamsFilters struct {
	UserID      uint   `json:"user_id"`      // Filtr dla wiadomosci przeczytanych przez uzytkownika
	Page        int    `json:"page"`         // Numer strony dla paginacji
	PageSize    int    `json:"page_size"`    // Ilość elementów na stronie
	OrderBy     string `json:"order_by"`     // Kryterium sortowania []
	StartDate   string `json:"start_date"`   // Filtr dla wiadomości wysłanych po danej dacie
	EndDate     string `json:"end_date"`     // Filtr dla wiadomości wysłanych przed daną datą
	ContentLike string `json:"content_like"` // Szukaj wiadomości zawierających podany ciąg znaków
	SenderID    uint   `json:"sender_id"`    // Filtruj wiadomości według ID nadawcy
}
