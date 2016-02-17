package main

import (
	"game"
)

type User struct {
	var userName string
	var password string
	var userId int
	var games []Game
	var messages []string
	var friends []User
	 

	func InitUser(user_name string, pass string) {
		userName = user_name
		password = pass
		//TODO: Figure out a static id
		userId = 1
	}

	func AddGame(game Game)	{
		games = append(games,game)
	}

	func AddFriend(friend User){
		friends = append(friends, friend)
	}

	func AddMessage(message String){
		messages = append(messages, message)
	}

}
