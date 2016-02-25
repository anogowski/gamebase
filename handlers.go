package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	tag := r.URL.Query().Get("tag")
	nameinclude := r.URL.Query().Get("NameIncludes")
	res := []models.Game{}
	var err error
	if tag!=""{
		//var tagres []models.Game
		//tagres, err = models.GlobalGameStore.FindTagged(tag)
		//res = append(res, tagres)
	}
	if nameinclude!=""{
		//var namres []models.Game
		//namres, err = models.GlobalGameStore.FindNameIncludes(nameinclude)
		//res = append(res, namres)
	}
	models.RenderTemplate(w, r, "game/search", map[string]interface{}{"NameIncludes":nameinclude, "Tag":tag, "SearchResults":res, "Error":err})
}

func HandleAccountPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w, r){
		models.RenderTemplate(w, r, "users/account", nil)
	}
}
func HandleAccountAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w, r){
		//TODO: handle account update
		
	}
}

func HandleChatAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		//TODO: send the chatNewMsg to chatTo
		
	}
}
