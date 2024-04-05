package app

import (
	"errors"
	"time"
)

type ChatroomList struct {
	Chatrooms []Chatroom
}

func (chatrooms *ChatroomList) addChatroom(chatroom Chatroom) {
	chatrooms.Chatrooms = append(chatrooms.Chatrooms, chatroom)
}

func (chatrooms *ChatroomList) getAllChatroomIDs() []int {
	var chatroomList []int
	for _, chatroom := range chatrooms.Chatrooms {
		chatroomList = append(chatroomList, chatroom.ID)
	}
	return chatroomList
}

func (chatrooms *ChatroomList) getChatroomByID(id int) (*Chatroom, error) {
	for i, chatroom := range chatrooms.Chatrooms {
		if chatroom.ID == id {
			return &chatrooms.Chatrooms[i], nil
		}
	}
	return nil, errors.New("could not fetch chatroom with given id")
}

type Chatroom struct {
	ID       int
	HostID   int
	Name     string
	Members  []*User
	Messages []Message
	Parent   *Chatroom
	Children []*Chatroom
}

func createChatroom(hostID int, members []*User, allChatrooms *ChatroomList) Chatroom {
	chatroomIDList := allChatrooms.getAllChatroomIDs()
	var chatroom Chatroom
	chatroom.ID = generateUniqueID(&chatroomIDList)
	chatroom.HostID = hostID
	chatroom.Members = append(chatroom.Members, members...)
	appendChatroomIDToUsers(chatroom.ID, members)
	return chatroom
}

func appendChatroomIDToUsers(chatroomID int, members []*User) {
	for _, member := range members {
		member.ChatroomIDs = append(member.ChatroomIDs, chatroomID)
	}
}

func (chatroom *Chatroom) addMessage(message Message) {
	message.Timestamp = time.Now()
	chatroom.Messages = append(chatroom.Messages, message)
}
