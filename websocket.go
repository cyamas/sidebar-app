package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	MsgType   string    `json:"msg-type"`
	RoomID    int       `json:"room-id"`
	UserID    int       `json:"user-id"`
	MemberIDs []int     `json:"member-ids"`
	Msg       []byte    `json:"message"`
	Received  time.Time `json:"timestamp"`
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket: ", err)
	}
	conn.WriteMessage(websocket.TextMessage, []byte("connected"))
	defer conn.Close()
	hub.Conns[conn] = true
	for {
		wsMessage, err := readAndUnmarshalMessage(conn)
		if err != nil {
			log.Println("Error reading or unmarshaling ws message: ", err)
			break
		}

		if wsMessage.MsgType == "text" {
			broadcastTextMessage(&wsMessage)
		}
		if wsMessage.MsgType == "connected" {
			user := hub.Users[wsMessage.UserID]
			user.WSConn = conn
			log.Println("User has connected")
		}
		if wsMessage.MsgType == "createroom" {
			createNewChatroom(&wsMessage, conn)
			host := hub.Users[wsMessage.UserID].Name
			log.Printf("%v has created a room!", host)
		}
	}
}

func readAndUnmarshalMessage(conn *websocket.Conn) (WSMessage, error) {
	var wsMessage WSMessage
	_, wsMsg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message: ", err)
		return wsMessage, err
	}
	if err := json.Unmarshal(wsMsg, &wsMessage); err != nil {
		log.Println("Could not unmarshal data: ", err)
		return wsMessage, err
	}
	return wsMessage, nil
}

func broadcastTextMessage(wsMessage *WSMessage) {
	room := hub.Chatrooms[wsMessage.RoomID]
	for _, user := range room.Members {
		if user.ID != wsMessage.UserID {
			conn := user.WSConn
			err := conn.WriteMessage(websocket.TextMessage, wsMessage.Msg)
			if err != nil {
				log.Println("Error writing message to websocket: ", err)
				return
			}
		}
	}
}

func createNewChatroom(wsMessage *WSMessage, conn *websocket.Conn) {
	newChatroom := hub.createChatroom(wsMessage.UserID, wsMessage.MemberIDs)
	data := make(map[string]int)
	data["roomID"] = newChatroom.ID
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Could not marshal data", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, jsonBytes)
}
