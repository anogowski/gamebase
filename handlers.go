package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	tag := r.URL.Query().Get("tag")
	//res, err := models.GlobalGameStore.FindTagged(tag)
	models.RenderTemplate(w, r, "game/search", map[string]interface{}{"Tag":tag/*, "SearchResults":res, "Error":err*/})
}
