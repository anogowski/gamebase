package main
import (
	
)

type SessionStore interface{
	Find(id string)(*Session, error)
	Delete(sess *Session)error
}
var globalSessionStore SessionStore

//TODO: implement the SessionStore interface and set the globalSessionStore
