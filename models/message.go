package models

type Message struct {
	To         User
	From       User
	TheMessage string
	Timestamp  Time
}

func NewMessage(to, from User, theMessage string) *Message {
	return &Message{To: to, From: from, TheMessage: theMessage}
}
