package main

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Conns     map[*websocket.Conn]bool
	Chatrooms map[int]*Chatroom
	Users     map[int]*User
}

func newHub() *Hub {
	return &Hub{
		Conns:     make(map[*websocket.Conn]bool),
		Chatrooms: make(map[int]*Chatroom),
		Users:     make(map[int]*User),
	}
}

func (hub *Hub) createUser(name string) *User {
	var user User
	allUserIDs := hub.getAllUserIDs()
	user.ID = generateUniqueID(allUserIDs)
	user.Name = name
	hub.Users[user.ID] = &user
	return &user
}

func (hub *Hub) getAllUserIDs() map[int]bool {
	userIDs := make(map[int]bool)
	for id := range hub.Users {
		userIDs[id] = true
	}
	return userIDs
}

func (hub *Hub) getAllUsernames() map[string]bool {
	usernames := make(map[string]bool)
	for _, user := range hub.Users {
		usernames[user.Name] = true
	}
	return usernames
}

func (hub *Hub) createChatroom(hostID int, memberIDs []int) *Chatroom {
	allChatroomIDs := hub.getAllChatroomIDs()
	var chatroom Chatroom
	chatroom.ID = generateUniqueID(allChatroomIDs)
	members := addMembersByID(hub.Users, memberIDs, chatroom.ID)
	chatroom.Host = hub.Users[hostID]
	chatroom.Members = append(chatroom.Members, members...)
	hub.Chatrooms[chatroom.ID] = &chatroom
	return &chatroom
}

func addMembersByID(allUsers map[int]*User, userIDs []int, chatroomID int) []*User {
	var members []*User
	for _, id := range userIDs {
		allUsers[id].ChatroomIDs = append(allUsers[id].ChatroomIDs, chatroomID)
		members = append(members, allUsers[id])
	}
	return members
}

func (hub *Hub) getAllChatroomIDs() map[int]bool {
	chatroomIDs := make(map[int]bool)
	for id := range hub.Chatrooms {
		chatroomIDs[id] = true
	}
	return chatroomIDs
}

func generateUniqueID(idList map[int]bool) int {
	rand.New(rand.NewSource(time.Now().Unix()))
	candidate := rand.Intn(900000) + 100000
	if IsUniqueID(candidate, idList) {
		return candidate
	} else {
		return generateUniqueID(idList)
	}
}

func IsUniqueID(candidate int, allIDs map[int]bool) bool {
	if len(allIDs) > 0 {
		_, ok := allIDs[candidate]
		if ok {
			return false
		}
	}
	return true
}
