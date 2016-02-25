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
	models.GlobalTagStore = models.NewPostgresTagStore()
}

func main() {
	router := httprouter.New()
	router.Handle("GET", "/login", HandleLoginPage)
	router.Handle("POST", "/login", HandleLoginAction)
	router.Handle("GET", "/logout", HandleLogout)
	router.Handle("GET", "/search", HandleSearch)
	router.Handle("GET", "/account", HandleAccountPage)
	router.Handle("POST", "/account", HandleAccountAction)
	router.Handle("POST", "/chat", HandleChatAction)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	fmt.Println("Server Running...")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
