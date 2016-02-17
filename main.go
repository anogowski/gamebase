package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

var PORT int
func init(){
	if PORT==0{
		PORT = 8080
	}
}

func main() {
	router := httprouter.New()
	router.Handle("GET", "/login", HandleLoginPage)
	router.Handle("POST", "/login", HandleLoginAction)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	
	fmt.Println("Server Running...")
	http.ListenAndServe(":"+string(PORT), router)
}
