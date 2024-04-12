package app

import "time"

type Message struct {
	SenderID  int
	Text      string
	Timestamp time.Time
}
