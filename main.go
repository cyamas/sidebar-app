package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	MsgType   string `json:"msg-type"`
	RoomID    int    `json:"room-id"`
	UserID    int    `json:"user-id"`
	MemberIDs []int  `json:"member-ids"`
	Msg       string `json:"message"`
}

type User struct {
	ID          int
	WSConn      *websocket.Conn
	Name        string
	ChatroomIDs []int
	Hub         *Hub
	Send        chan ClientMessage
}

func (user *User) ReadPump() {

	for {
		_, wsMsg, err := user.WSConn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading message: ", err)
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(wsMsg, &clientMsg); err != nil {
			log.Println("Could not unmarshal data: ", err)
		}
		msgType := clientMsg.MsgType
		if msgType == "text" {
			room, err := user.Hub.getChatroomByID(clientMsg.RoomID)
			if err != nil {
				log.Println(err)
			}
			room.Broadcast <- clientMsg
		}
	}
}

func (user *User) WritePump() {

}

type Chatroom struct {
	ID        int
	Host      *User
	Name      string
	Members   []*User
	Messages  []TextMessage
	Parent    *Chatroom
	Children  []*Chatroom
	Broadcast chan ClientMessage
}

type TextMessage struct {
	MsgType    string
	Msg        string
	SenderID   int
	SenderName string
	RoomID     int
}

var hub = newHub()

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Println("Sidebar Hub has started.")
	router := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("static"))
	router.Use(middleware.Logger)
	router.Get("/", home)
	router.Post("/signin", signinUser)
	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(w, r)
	})
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	http.ListenAndServe(":6699", router)
}
