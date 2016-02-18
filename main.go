package main

import (
	"os"
	"log"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

var PORT string
func init(){
	PORT = os.Getenv("PORT")
	if PORT==""{
		PORT = "8080"
	}
}

func main() {
	router := httprouter.New()
	router.Handle("GET", "/login", HandleLoginPage)
	router.Handle("POST", "/login", HandleLoginAction)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	
	fmt.Println("Server Running...")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
