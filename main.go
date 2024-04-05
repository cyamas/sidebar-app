package main

import (
	"fmt"
	"log"

	"github.com/sidebar-app/app"
)

func main() {
	var allUsers app.AllUsers
	var allChatrooms app.ChatroomList

	userIDList := allUsers.GetAllUserIDs()
	allUsers.CreateUser("shgi24", &userIDList)
	allUsers.CreateUser("Waffleman26", &userIDList)
	allUsers.CreateUser("Leoguy14", &userIDList)
	user1, err := allUsers.GetUserByID(allUsers.Users[0].ID)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	user2, err := allUsers.GetUserByID(allUsers.Users[1].ID)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	user3, err := allUsers.GetUserByID(allUsers.Users[2].ID)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	allMembers := []*app.User{user1, user2, user3}
	user1.HostNewChatroom(allMembers, &allChatrooms)
	chatroom1 := &allChatrooms.Chatrooms[0]
	chatroomOneID := chatroom1.ID
	user2.SendMessage("Hello, friends!", chatroomOneID, &allChatrooms)
	user3.SendMessage("AHHHHHHWOOAHHAA!", chatroomOneID, &allChatrooms)
	user1.SendMessage("Do either of you know why birds fly gladly in a V?", chatroomOneID, &allChatrooms)
	user2.HostSubChatroom(chatroomOneID, []*app.User{user2, user3}, &allChatrooms)
	chatroom2 := &allChatrooms.Chatrooms[1]
	user2.SendMessage("Your dad is crazy, Leoguy...", chatroom2.ID, &allChatrooms)
	user3.SendMessage("AHHGAAA", chatroom2.ID, &allChatrooms)

	fmt.Printf("All allUsers: %v\n", allUsers)
	fmt.Println("\nchatroom1 messages:")
	for _, message := range chatroom1.Messages {
		fmt.Println(message)
	}
	fmt.Println("\nchatroom2 messages:")
	for _, message := range chatroom2.Messages {
		fmt.Println(message)
	}
}
