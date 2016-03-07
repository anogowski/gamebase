package main

import (
	"fmt"
	"github.com/anogowski/gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/anogowski/gamebase/models"
	"log"
	"net/http"
	"os"
)

var PORT string

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	models.GlobalUserStore = models.NewPostgresUserStore()
	models.GlobalSessionStore = models.NewPostgresSessionStore()
}

func main() {
	router := httprouter.New()
	router.Handle("GET", "/", HandleIndex)
	router.Handle("GET", "/game/:wild", HandleGamePage)
	router.Handle("POST", "/game/:wild", HandleGamePageNewAction)
	router.Handle("GET", "/game/:wild/edit", HandleGameEditPage)
	router.Handle("POST", "/game/:wild/edit", HandleGameEditAction)
	router.Handle("GET", "/review/:wild/new", HandleReviewNew)
	router.Handle("POST", "/review/:wild/new", HandleReviewNewAction)
	router.Handle("GET", "/review/:wild", HandleReview)
	router.Handle("GET", "/videos/:wild/new", HandleVideoNew)
	router.Handle("POST", "/videos/:wild/new", HandleVideoNewAction)
	router.Handle("GET", "/videos/:wild", HandleVideo)
	router.Handle("GET", "/login", HandleLoginPage)
	router.Handle("POST", "/login", HandleLoginAction)
	router.Handle("GET", "/logout", HandleLogout)
	router.Handle("GET", "/search", HandleSearch)
	router.Handle("GET", "/account", HandleAccountPage)
	router.Handle("POST", "/account", HandleAccountAction)
	router.Handle("POST", "/chat", HandleChatAction)
	router.Handle("GET", "/users/:wild", HandleUserPage)
	router.Handle("GET", "/friend/:wild", HandleFriendAdd)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	fmt.Println("Server Running...")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
