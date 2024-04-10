package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/sidebar-app/app"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = newHub()

func main() {
	router := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("static"))
	router.Use(middleware.Logger)
	router.Get("/", home)
	router.Get("/ws", handleConnection)
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	http.ListenAndServe(":6699", router)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket: ", err)
	}
	greeting := []byte("Hello from the sidebar server")
	err = conn.WriteMessage(websocket.TextMessage, greeting)
	if err != nil {
		log.Println("Could not send greeting: ", err)
	}
	defer conn.Close()
	hub.Conns[conn] = true
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		log.Println(messageType, string(p))
	}
}

type Hub struct {
	Conns     map[*websocket.Conn]bool
	Chatrooms map[*app.Chatroom]bool
}

func newHub() *Hub {
	return &Hub{
		Conns:     make(map[*websocket.Conn]bool),
		Chatrooms: make(map[*app.Chatroom]bool),
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := "OK"
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
