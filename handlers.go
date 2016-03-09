package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"net/url"
)

func HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	randgames := []models.Game{}
	var err error
	//randgames,err = models.Dal.GetGames(20)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w, r, "home/index", map[string]interface{}{"RandomGames":randgames})
}

func HandleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	tag := r.URL.Query().Get("tag")
	nameinclude := r.URL.Query().Get("NameIncludes")
	res := []models.Game{}
	var err error
	if tag!=""{
		var tagres []models.Game
		tagres, err = models.Dal.FindGamesByTag(tag)
		res = append(res, tagres...)
	}
	if nameinclude!=""{
		//var namres []models.Game
		//namres, err = models.Dal.SearchGames(nameinclude)
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
		//update user information
		username := r.FormValue("accountName")
		email := r.FormValue("accountEmail")
		newPassword := r.FormValue("accountNewPassword")
		confirmNewPassword := r.FormValue("confirmPassword")
		oldPassword := r.FormValue("accountPassword")

		
		user, err := models.GlobalUserStore.Authenticate(username, oldPassword)
		if err != nil {
			models.RenderTemplate(w, r, "users/account", map[string]interface{}{"Error": err.Error()})
			return
		}
		
		user.Email = email
		if newPassword != ""{
			if(newPassword != confirmNewPassword){
				models.RenderTemplate(w, r, "users/account", map[string]interface{}{"Error": "Passwords do not match."})
				return
			}
			user.SetPassword(newPassword)
		}
		
		err = models.Dal.UpdateUser(*user)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w, r, "users/account", nil)
	}
}

func HandleChatAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		//TODO: send the chatNewMsg to chatTo
		user := models.RequestUser(r)
		toUser := r.FormValue("chatTo")
		chatID := r.FormValue("chatToID")
		theMessage := r.FormValue("chatNewMsg")
		models.Dal.SendMessage(user, chatID, theMessage)

	}
}

func HandleGamePage(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	encgameid := params.ByName("wild")
	if encgameid=="new"{
		HandleGamePageNew(w,r,params)
		return
	}
	var err error
	gameid, err := url.QueryUnescape(encgameid)
	if err!=nil{
		panic(err)
	}
	var game *models.Game
	game, err = models.Dal.FindGame(gameid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r,"game/page", map[string]interface{}{"Game":game})
}
func HandleGamePageNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		models.RenderTemplate(w,r, "game/new", nil)
	}
}
func HandleGamePageNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if params.ByName("wild")=="new" && models.SignedIn(w,r){
		title := r.FormValue("gameTitle")
		dev := r.FormValue("gameDeveloper")
		pub := r.FormValue("gamePublisher")
		trailer := r.FormValue("gameTrailer")
		//copy := r.FormValue("gameCopyright")
		desc := r.FormValue("gameDescription")
		tagstr := r.FormValue("gameTags")
		var tags []string
		err := json.Unmarshal([]byte(tagstr), &tags)
		if err!=nil || tags==nil{
			tags = []string{}
		}
		game := models.NewGame(title, pub, trailer)
		game.Description = desc
		game.Developer = dev
		err = models.Dal.CreateGame(*game)
		if err!=nil{
			panic(err)
		}
		for _,tag := range tags{
			models.Dal.AddGameTag(game.GameId, tag)
		}
		http.Redirect(w,r, "/game/"+url.QueryEscape(game.GameId), http.StatusFound)
	} else{
		http.NotFound(w,r)
	}
}
func HandleGameEditPage(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		game, err := models.Dal.FindGame(gameid)
		if err!=nil{
			panic(err)
		}
		if game==nil{
			http.Redirect(w,r, "/?flash=Game+Not+Found", http.StatusNotFound)
			return
		}
		tags := []string{}
		tags, err = models.Dal.FindTagsByGame(gameid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "game/edit", map[string]interface{}{"Game":game, "Tags":tags})
	}
}
func HandleGameEditAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		title := r.FormValue("gameTitle")
		dev := r.FormValue("gameDeveloper")
		pub := r.FormValue("gamePublisher")
		trailer := r.FormValue("gameTrailer")
		//copy := r.FormValue("gameCopyright")
		desc := r.FormValue("gameDescription")
		tagstr := r.FormValue("gameTags")
		var tags []string
		err := json.Unmarshal([]byte(tagstr), &tags)
		if err!=nil || tags==nil{
			tags = []string{}
		}
		for k,v := range tags{
			t,err := url.QueryUnescape(v)
			if err==nil{
				tags[k] = t
			}
		}
		var currtags []string
		currtags, err = models.Dal.FindTagsByGame(gameid)
		if err!=nil || currtags==nil{
			currtags = []string{}
		}
		game := models.Game{GameId:gameid, Title:title, Developer:dev, Publisher:pub, URL:trailer, Description:desc}
		err = models.Dal.UpdateGame(game)
		if err!=nil{
			panic(err)
		}
		for _,t := range currtags{
			found := false
			for k,v := range tags {
				if t==v{
					found = true
					tags = append(tags[:k], tags[k+1:]...)
					break;
				}
			}
			if !found{
				models.Dal.DeleteGameTag(gameid, t)
			}
		}
		for _,t := range tags{
			models.Dal.AddGameTag(gameid, t)
		}
		http.Redirect(w,r, "/game/"+url.QueryEscape(game.GameId), http.StatusFound)
	}
}
func HandleReview(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	//reviewid := params.ByName("wild")
	//rev, err := models.GlobalReviewStore.Find(reviewid)
	//models.RenderTemplate(w,r, "review/review", map[string]interface{}{"Review":rev})
}
func HandleReviewNew(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		var gam *models.Game
		var err error
		//gam, err := models.GlobalGameStore.Find(gameid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "review/new", map[string]interface{}{"Game":gam})
	}
}
func HandleReviewNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: create the new review page
		//reviewid
		reviewid := r.FormValue("reviewID")
		//userid
		user := models.RequestUser(r)
		//gameid
		gameid := r.FormValue("gameID")
		//review
		review := r.FormValue("reviewBody")
		//rating
		rating := r.FormValue("reviewRating")
		//models.Dal.CreateReview()
	}
}
func HandleVideo(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	videoid := params.ByName("wild")
	user, err := models.GlobalUserStore.FindUser(videoid)
	if err!=nil{
		panic(err)
	}
	if user!=nil{
		var vids []models.Video
		//vids, err = models.GlobalVideoStore.FindByUser(user.UserId)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "videos/alluservideos", map[string]interface{}{"User":user, "AllVideos":vids})
	} else{
		var vid *models.Video
		var err error
		//vid, err = models.GlobalVideoStore.Find(videoid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "video/uservideo", map[string]interface{}{"Video":vid, "VideoWidth":640, "VideoHeight":480})
	}
}
func HandleVideoNew(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		var gam *models.Game
		var err error
		//game, err = models.GlobalGameStore.Find(gameid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "video/new", map[string]interface{}{"Game":gam})
	}
}
func HandleVideoNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		//gameid := params.ByName("wild")
		//TODO: create the new video page
		
	}
}

func HandleUserPage(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	userid := params.ByName("wild")
	user, err := models.GlobalUserStore.FindUser(userid)
	if err!=nil{
		panic(err)
	}
	var vids []models.Video
	var revs []models.Review
	//vids, err = models.GlobalVideoStore.FindByUser(userid)
	if err!=nil{
		panic(err)
	}
	//revs, err = models.GlobalReviewStore.FindByUser(userid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r, "users/view", map[string]interface{}{"User":user, "Reviews":revs, "Videos":vids})
}
func HandleFriendAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		friendid := params.ByName("wild")
		friend, err := models.GlobalUserStore.FindUser(friendid)
		if err!=nil{
			panic(err)
		}
		user := models.RequestUser(r)
		if user==nil || friend==nil || user.UserId==friend.UserId{
			models.RenderTemplate(w,r, "users/friend", nil)
			return
		}
		//err = models.FriendStore.AddFriend(user.UserId, friendid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "users/friend", map[string]interface{}{"Friend":friend})
	}
}
