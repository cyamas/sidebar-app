package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type User struct {
	ID          int
	WSConn      *websocket.Conn
	Name        string
	ChatroomIDs []int
}

type Chatroom struct {
	ID       int
	Host     *User
	Name     string
	Members  []*User
	Messages []Message
	Parent   *Chatroom
	Children []*Chatroom
}

type Message struct {
	SenderID int
	Text     string
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
