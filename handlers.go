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
		//update user information
		username := r.FormValue("accountName")
		email := r.FormValue("email")
		newPassword := r.FormValue("accountNewPassword")
		oldPassword := r.FormValue("accountPassword")

		user, err = models.GlobalUserStore.Authenticate(username, oldPassword)

		if err != nil {
			models.RenderTemplate(w, r, "users/account", map[string]interface{}{"Error": err.Error()})
			return
		}

		user.Email = email
		if newPassword != ""{
			hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), hashcost)
			newPassword = string(hash)
			user.SetPassword(newPassword)
		}

		err := models.Dal.UpdateUser(user)
		models.RenderTemplate(w, r, "users/account", nil)
		//call the dal for update next
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
	var game *models.Game
	//game, err := models.GlobalGameStore.Find(gameid)
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
func HandleGamePageNewAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		title := r.FormValue("gameTitle")
		pub := r.FormValue("gamePublisher")
		trailer := r.FormValue("gameTrailer")
		//copy := r.FormValue("gameCopyright")
		//desc := r.FormValue("gameDescription")
		tagstr := r.FormValue("gameTags")
		var tags []string
		err := json.Unmarshal([]byte(tagstr), &tags)
		if err!=nil || tags==nil{
			tags = []string{}
		}
		game := models.NewGame(title, pub, trailer)
		//models.DAL.CreateGame(game)
		//for _,tag := range tags{
		//	models.DAL.AddGameTag(game.GameId, tag)
		//}
		http.Redirect(w,r, "/game/"+url.QueryEscape(game.GameId), http.StatusFound)
	}
}
func HandleGameEditPage(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		game, err := models.DAL.FindGame(gameid)
		if err!=nil{
			panic(err)
		}
		if game==nil{
			http.Redirect(w,r, "/?flash=Game+Not+Found", http.StatusNotFound)
			return
		}
		tags := []string{}
		//tags, err := models.DAL.FindTagsByGame(gameid)
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
		pub := r.FormValue("gamePublisher")
		trailer := r.FormValue("gameTrailer")
		//copy := r.FormValue("gameCopyright")
		//desc := r.FormValue("gameDescription")
		tagstr := r.FormValue("gameTags")
		var tags []string
		err := json.Unmarshal([]byte(tagstr), &tags)
		if err!=nil || tags==nil{
			tags = []string{}
		}
		remtagstr := r.FormValue("removedTags")
		var remtags []string
		err = json.Unmarshal([]byte(remtagstr), &tags)
		if err!=nil || remtags==nil{
			remtags = []string{}
		}
		game := models.NewGame(title, pub, trailer)
		err = models.DAL.UpdateGame(game)
		if err!=nil{
			panic(err)
		}
		//for _,t := range remtags{
		//	models.DAL.RemoveGameTag(gameid, t)
		//}
		//for _,t := range tags{
		//	models.DAL.AddGameTag(gameid, t)
		//}
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
