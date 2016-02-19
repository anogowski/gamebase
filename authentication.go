package main

import (
	"github.com/anogowski/gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/anogowski/gamebase/models"
	"log"
	"net/http"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	models.RenderTemplate(w, r, "users/login", nil)
}
func HandleLoginAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}
	if r.URL.Query().Get("signup") == "true" {
		uname := r.FormValue("signUser")
		pword := r.FormValue("signPass")
		repeat := r.FormValue("repeatPass")
		if pword != repeat {
			models.RenderTemplate(w, r, "/users/login", map[string]interface{}{"Error": "Passwords don't match.", "UName": uname})
			return
		}
		user, err := models.GlobalUserStore.FindUserByName(uname)
		if err != nil {
			panic(err)
		}
		if user != nil {
			models.RenderTemplate(w, r, "/users/login", map[string]interface{}{"Error": "Username not available.", "UName": uname})
			return
		}
		user, err = models.GlobalUserStore.CreateUser(uname, pword)
		if err != nil {
			panic(err)
		}
		if user != nil {
			log.Fatal("Failed to create user.")
		}
		http.Redirect(w, r, next+"?flash=Signup+Success", http.StatusFound)
	} else {
		uname := r.FormValue("loginUser")
		pword := r.FormValue("loginPass")
		_, err := models.GlobalUserStore.Authenticate(uname, pword)
		if err != nil {
			models.RenderTemplate(w, r, "/users/login", map[string]interface{}{"Error": err.Error(), "UName": uname})
			return
		}
		http.Redirect(w, r, next+"?flash=Login+Success", http.StatusFound)
	}
}
