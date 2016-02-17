package main
import (
	
)

type UserStore interface{
	CreateUser(name, pass string)(*User, error)
	FindUser(id string)(*User,error)
	FindUserByName(name string)(*User, error)
	Authenticate(name, pass string)(*User, error)
}

var globalUserStore UserStore

//TODO: implement the UserStore interface and set the globalUserStore
