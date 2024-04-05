package app

import "time"

type Message struct {
	SenderID  int
	Text      string
	Timestamp time.Time
}

func createMessage(text string, userID int) Message {
	var message Message
	message.SenderID = userID
	message.Text = text
	message.Timestamp = time.Now()
	return message
}
