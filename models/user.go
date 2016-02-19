package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName string
	Password string
	UserId string
	Email string
	Games []Game
	Messages []string
	Friends []User
}

const (
	hashcost = 10
	userIDlen = 20
)

func NewUser(user_name, pass string) *User{
	user := User{}
	user.InitUser(user_name, pass)
	return &user
}
func (this *User) InitUser(user_name string, pass string) {
	this.UserName = user_name
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), hashcost)
	this.Password = string(hash)
	this.UserId = GenerateID("user_", userIDlen)
}

func (this *User) AddGame(game Game)	{
	this.Games = append(this.Games,game)
}

func (this *User) AddFriend(friend User){
	this.Friends = append(this.Friends, friend)
}

func (this *User) AddMessage(message string){
	this.Messages = append(this.Messages, message)
}
func (this *User) CheckPassword(pass string)bool{
	return bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(pass))==nil
}
