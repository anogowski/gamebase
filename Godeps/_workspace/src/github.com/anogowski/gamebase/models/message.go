package models

type Message struct{
	To User
	From User
	TheMessage string
}

func (this *Message) NewMessage(to, from User, theMessage string){
	return message{To: to, From: from, TheMessage: theMessage}
}