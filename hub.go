package main

import (
	"errors"
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
		Chatrooms: make(map[int]*Chatroom),
		Users:     make(map[int]*User),
	}
}

func (hub *Hub) createUser(name string) *User {
	newUserID := generateUniqueID(hub.getAllUserIDs())
	user := User{
		ID:   newUserID,
		Name: name,
		Hub:  hub,
		Send: make(chan ClientMessage),
	}
	hub.Users[newUserID] = &user
	return &user
}

func (hub *Hub) getUserByID(id int) (*User, error) {
	if user, ok := hub.Users[id]; ok {
		return user, nil
	}
	err := errors.New("no user associated with id")
	return nil, err
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
	newRoomID := generateUniqueID(hub.getAllChatroomIDs())
	chatroom := Chatroom{
		ID:        newRoomID,
		Host:      hub.Users[hostID],
		Members:   addMembersByID(hub.Users, memberIDs, newRoomID),
		Broadcast: make(chan ClientMessage),
	}

	members := addMembersByID(hub.Users, memberIDs, chatroom.ID)
	chatroom.Members = append(chatroom.Members, members...)

	hub.Chatrooms[chatroom.ID] = &chatroom
	return &chatroom
}

func (hub *Hub) getChatroomByID(id int) (*Chatroom, error) {
	room, ok := hub.Chatrooms[id]
	if ok {
		return room, nil
	}
	err := errors.New("could not find chatroom with given id")
	return nil, err
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
