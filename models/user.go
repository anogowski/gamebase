package main

import (
	
)

type User struct {
	userName string
	password string
	userId string
	games []Game
	messages []string
	friends []User
}

func (this *User) InitUser(user_name string, pass string) {
	this.userName = user_name
	this.password = pass
	//TODO: Figure out a static id, maybe GenerateID("user_", 20)
	this.userId = string(1)
}

func (this *User) AddGame(game Game)	{
	this.games = append(this.games,game)
}

func (this *User) AddFriend(friend User){
	this.friends = append(this.friends, friend)
}

func (this *User) AddMessage(message string){
	this.messages = append(this.messages, message)
}
