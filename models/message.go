package models
import (
	"time"
)

type Message struct {
	To         User
	From       User
	TheMessage string
	TimeStamp  time.Time
}

func NewMessage(to, from User, theMessage string) *Message {
	return &Message{To: to, From: from, TheMessage: theMessage, TimeStamp:time.Now()}
}
