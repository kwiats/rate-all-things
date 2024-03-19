package model

import (
	"nhooyr.io/websocket"
)

type Message struct {
	Reciver Reciver `json:"reciver"`
	Sender  int     `json:"sender"`
	Content string  `json:"content"`
}

type Reciver struct {
	Type string `json:"type"` // group/user
	Id   uint   `json:"id"`
}

type Client struct {
	Id         uint
	Connection *websocket.Conn
}

type Group struct {
	Id      string
	Members map[*Client]bool
}
