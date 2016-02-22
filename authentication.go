package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	models.RenderTemplate(w, r, "users/login", nil)
}
func HandleLoginAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user *models.User
	flash := ""
	if r.URL.Query().Get("signup") == "true" {
		uname := r.FormValue("signUser")
		email := r.FormValue("signEmail")
		pword := r.FormValue("signPass")
		repeat := r.FormValue("repeatPass")
		if pword != repeat {
			models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": "Passwords don't match.", "UName": uname})
			return
		}
		user, err := models.GlobalUserStore.FindUserByName(uname)
		if err != nil {
			panic(err)
		}
		if user != nil {
			models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": "Username not available.", "UName": uname})
			return
		}
		user, err = models.GlobalUserStore.CreateUser(uname, pword, email)
		if err != nil {
			panic(err)
		}
		if user != nil {
			log.Fatal("Failed to create user.")
		}
		flash = "?flash=Signup+Success"
	} else {
		uname := r.FormValue("loginUser")
		pword := r.FormValue("loginPass")
		user, err := models.GlobalUserStore.Authenticate(uname, pword)
		if err != nil {
			models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": err.Error(), "UName": uname})
			return
		}
		flash = "?flash=Login+Success"
	}
	sess := models.FindOrCreateSession(w, r)
	sess.UserID = user.UserId
	err := models.GlobalSessionStore.Save(sess)
	if err!=nil{
		panic(err)
	}
	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+flash, http.StatusFound)
}
func HandleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	sess := models.RequestSession(r)
	if sess!=nil{
		err := models.GlobalSessionStore.Delete(sess)
		if err!=nil{
			panic(err)
		}
	}
	models.RenderTemplate(w,r, "users/logout", nil)
}
