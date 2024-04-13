package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	MsgType   string `json:"msg-type"`
	RoomID    int    `json:"room-id"`
	UserID    int    `json:"user-id"`
	MemberIDs []int  `json:"member-ids"`
	Msg       string `json:"message"`
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn := connectWS(w, r)
	defer conn.Close()
	hub.Conns[conn] = true
	confirmConnection(conn)
	for {
		clientMsg, err := readAndUnmarshalMessage(conn)
		if err != nil {
			log.Println("Error reading or unmarshaling ws message: ", err)
			break
		}
		handleClientMessage(&clientMsg, conn)
	}
}

func connectWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket: ", err)
	}
	hub.Conns[conn] = true
	return conn
}

func confirmConnection(conn *websocket.Conn) {
	connectedMsg := map[string]string{"MsgType": "connected"}
	marshaledMsg, err := json.Marshal(connectedMsg)
	if err != nil {
		log.Println("Error marshaling connectedMsg to json: ", err)
	}
	conn.WriteMessage(websocket.TextMessage, marshaledMsg)
}

func readAndUnmarshalMessage(conn *websocket.Conn) (ClientMessage, error) {
	var clientMsg ClientMessage
	_, wsMsg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message: ", err)
		return clientMsg, err
	}
	if err := json.Unmarshal(wsMsg, &clientMsg); err != nil {
		log.Println("Could not unmarshal data: ", err)
		return clientMsg, err
	}
	return clientMsg, nil
}

func handleClientMessage(clientMsg *ClientMessage, conn *websocket.Conn) {

	switch clientMsg.MsgType {
	case "connected":
		addConnToUser(clientMsg, conn)

	case "text":
		broadcastTextMessage(clientMsg)

	case "createroom":
		createNewChatroom(clientMsg, conn)
	}
}

func addConnToUser(clientMsg *ClientMessage, conn *websocket.Conn) {
	user := hub.Users[clientMsg.UserID]
	user.WSConn = conn
}

func broadcastTextMessage(clientMsg *ClientMessage) {
	room := hub.Chatrooms[clientMsg.RoomID]
	for _, user := range room.Members {
		if user.ID != clientMsg.UserID {
			conn := user.WSConn
			err := conn.WriteMessage(websocket.TextMessage, []byte(clientMsg.Msg))
			if err != nil {
				log.Println("Error writing message to websocket: ", err)
				return
			}
		}
	}
}

func createNewChatroom(clientMsg *ClientMessage, conn *websocket.Conn) {
	chatroom := hub.createChatroom(clientMsg.UserID, clientMsg.MemberIDs)
	newRoomMsg := createNewChatroomWSMessage(chatroom)
	marshaledNewRoomMsg, err := json.Marshal(newRoomMsg)
	if err != nil {
		log.Println("Could not marshal data", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, marshaledNewRoomMsg)
}

type WSNewChatroom struct {
	MsgType  string
	RoomID   int
	HostID   int
	HostName string
	Members  []map[int]string
}

func createNewChatroomWSMessage(chatroom *Chatroom) *WSNewChatroom {
	var newRoomMsg = WSNewChatroom{
		MsgType:  "newroom",
		HostID:   chatroom.Host.ID,
		HostName: chatroom.Host.Name,
	}
	newRoomMsg.RoomID = chatroom.ID
	var members []map[int]string
	for _, user := range chatroom.Members {
		member := make(map[int]string)
		member[user.ID] = user.Name
		members = append(members, member)
	}
	newRoomMsg.Members = members
	return &newRoomMsg
}
