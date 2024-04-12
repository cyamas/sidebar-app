package app

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID          int
	WSConn      *websocket.Conn
	Name        string
	ChatroomIDs []int
}
