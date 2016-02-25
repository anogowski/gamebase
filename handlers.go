package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"net/http"
	_"net/url"
)

func HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	randgames := []models.Game{}
	//randgames := models.GlobalGameStore.GetRandomGames(10)
	models.RenderTemplate(w, r, "home/index", map[string]interface{}{"RandomGames":randgames})
}

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

func HandleGamePage(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	encgameid := params.ByName("wild")
	if encgameid=="new"{
		HandleGamePageNew(w,r,params)
		return
	}
	var err error
	//gameid, err := url.QueryUnescape(encgameid)
	if err!=nil{
		panic(err)
	}
	var game models.Game
	//game, err := models.GlobalGameStore.Find(gameid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r,"game/page", map[string]interface{}{"Game":game})
}
func HandleGamePageNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		//TODO: handle showing the new game page
		
	}
}
func HandleGamePageNewAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		//TODO: handle creating the new game page
		
	}
}
func HandleReview(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	//reviewid := params.ByName("wild")
	//TODO: display the given review
	
}
func HandleReviewNew(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: display the new review page
		
	}
}
func HandleReviewNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: create the new review page
		
	}
}
func HandleVideo(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	//videoid := params.ByName("wild")
	//TODO: display the video page
	
}
func HandleVideoNew(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: display the new video page
		
	}
}
func HandleVideoNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: create the new video page
		
	}
}
