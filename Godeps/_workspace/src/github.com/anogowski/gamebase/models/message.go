package models

type Message struct{
	To User
	From User
	TheMessage string
}

func NewMessage(to, from User, theMessage string)*Message{
	return &message{To: to, From: from, TheMessage: theMessage}
}