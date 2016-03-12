package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/anogowski/gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/anogowski/gamebase/models"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Next": r.URL.Query().Get("next")})
}
func HandleLoginAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user *models.User
	next := string(r.FormValue("next"))
	flash := ""
	if r.URL.Query().Get("signup") == "true" {
		uname := r.FormValue("signUser")
		email := r.FormValue("signEmail")
		pword := r.FormValue("signPass")
		repeat := r.FormValue("repeatPass")
		if pword != repeat {
			models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": "Passwords don't match.", "UName": uname, "Email": email, "Next": next})
			return
		}
		var err error
		user, err = models.GlobalUserStore.CreateUser(uname, pword, email)
		if err != nil {
			if err==models.ErrUserNameUnavailable{
				models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": err.Error(), "UName": uname, "Email": email, "Next": next})
				return
			}
			panic(err)
		}
		if user == nil {
			log.Fatal("Failed to create user.")
		}
		flash = "?flash=Signup+Success"
	} else {
		uname := r.FormValue("loginUser")
		pword := r.FormValue("loginPass")
		var err error
		user, err = models.GlobalUserStore.Authenticate(uname, pword)
		if err != nil {
			models.RenderTemplate(w, r, "users/login", map[string]interface{}{"Error": err.Error(), "UName": uname, "Next": next})
			return
		}
		flash = "?flash=Login+Success"
	}
	next, err := url.QueryUnescape(next)
	if err != nil {
		next = ""
	}
	if next == "" {
		next = "/"
	}
	models.FindOrCreateSession(w, r, user.UserId)
	http.Redirect(w, r, next+flash, http.StatusFound)
}
func HandleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := models.RequestSession(r)
	if sess != nil {
		err := models.GlobalSessionStore.Delete(sess)
		if err != nil {
			panic(err)
		}
	}
	models.RenderTemplate(w, r, "users/logout", nil)
}
