package app

import (
	"errors"
	"log"
)

type AllUsers struct {
	Users []User
}

func (users *AllUsers) CreateUser(name string, userIDList *[]int) {
	var user User
	user.Name = name
	user.ID = generateUniqueID(userIDList)
	*userIDList = append(*userIDList, user.ID)
	users.Users = append(users.Users, user)
}

func (users *AllUsers) GetAllUserIDs() []int {
	var userIDList []int
	for _, user := range users.Users {
		userIDList = append(userIDList, user.ID)
	}
	return userIDList
}

func (users *AllUsers) GetUserByID(id int) (*User, error) {
	for i, user := range users.Users {
		if user.ID == id {
			return &users.Users[i], nil
		}
	}
	return nil, errors.New("could not fetch User from given id")
}

type User struct {
	ID          int
	Name        string
	ChatroomIDs []int
}

func (user *User) HostNewChatroom(members []*User, allChatrooms *ChatroomList) {
	chatroom := createChatroom(user.ID, members, allChatrooms)
	allChatrooms.addChatroom(chatroom)
}

func (user *User) HostSubChatroom(parentID int, members []*User, allChatrooms *ChatroomList) {
	parentChatroom, err := allChatrooms.getChatroomByID(parentID)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	subChatroom := createChatroom(user.ID, members, allChatrooms)
	parentChatroom.Children = append(parentChatroom.Children, &subChatroom)
	subChatroom.Parent = parentChatroom
	allChatrooms.addChatroom(subChatroom)
}

func (user *User) SendMessage(text string, chatroomID int, allChatrooms *ChatroomList) {
	chatroom, err := allChatrooms.getChatroomByID(chatroomID)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	message := createMessage(text, user.ID)
	chatroom.addMessage(message)
}
