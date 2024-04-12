package app

type Chatroom struct {
	ID       int
	Host     *User
	Name     string
	Members  []*User
	Messages []Message
	Parent   *Chatroom
	Children []*Chatroom
}
