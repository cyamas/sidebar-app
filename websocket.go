package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	confirmMsg := map[string]string{"MsgType": "connected"}
	marshaledMsg, err := json.Marshal(confirmMsg)
	if err != nil {
		log.Println("error marshaling map to json", err)
	}
	conn.WriteMessage(websocket.TextMessage, marshaledMsg)
	log.Println("Confirming connection...")
}

func readAndUnmarshalMessage(conn *websocket.Conn) (ClientMessage, error) {
	var clientMsg ClientMessage
	_, wsMsg, err := conn.ReadMessage()
	if err != nil {
		log.Fatalln("Error reading message: ", err)
		return clientMsg, err
	}
	if err := json.Unmarshal(wsMsg, &clientMsg); err != nil {
		log.Println("Could not unmarshal data: ", err)
		return clientMsg, err
	}
	log.Println("client msg: ", clientMsg)
	return clientMsg, nil
}

func handleClientMessage(clientMsg *ClientMessage, conn *websocket.Conn) {

	switch clientMsg.MsgType {
	case "connected":
		addConnToUser(clientMsg, conn)
		log.Println("Connection confirmed. Info added to corresponding user")

	case "text":
		broadcastTextMessage(clientMsg)
		log.Println(clientMsg.Msg)

	case "createroom":
		createNewChatroom(clientMsg, conn)

	case "activeusers":
		getActiveUsers(conn)
	}
}

func addConnToUser(clientMsg *ClientMessage, conn *websocket.Conn) {
	user := hub.Users[clientMsg.UserID]
	user.WSConn = conn
}

func broadcastTextMessage(clientMsg *ClientMessage) {
	msg := TextMessage{
		MsgType:  "text",
		Msg:      clientMsg.Msg,
		SenderID: clientMsg.UserID,
		RoomID:   clientMsg.RoomID,
	}
	user, err := hub.getUserByID(clientMsg.UserID)
	if err != nil {
		log.Println("User not found for given id.")
	}
	msg.SenderName = user.Name
	marshaledMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("could not marshal TextMessage for broadcast: ", err)
	}
	room := hub.Chatrooms[msg.RoomID]
	for _, user := range room.Members {
		if user.ID != msg.SenderID {
			err := user.WSConn.WriteMessage(websocket.TextMessage, []byte(marshaledMsg))
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

func getActiveUsers(conn *websocket.Conn) {
	allUsers := make(map[int]string)
	for id, user := range hub.Users {
		allUsers[id] = user.Name
	}
	msg := struct {
		MsgType string
		Users   map[int]string
	}{
		MsgType: "allusers",
		Users:   allUsers,
	}
	marshaledMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshaling data", err)
	}
	conn.WriteMessage(websocket.TextMessage, marshaledMsg)
}
