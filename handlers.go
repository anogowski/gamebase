package main

import (
	"gamebase/Godeps/_workspace/src/github.com/anogowski/gamebase/models"
	"gamebase/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	randgames := []models.Game{}
	var err error
	randgames,err = models.Dal.GetGames(20, 0)
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
		var namres []models.Game
		namres, err = models.Dal.SearchGames(nameinclude)
		res = append(res, namres...)
	}
	models.RenderTemplate(w, r, "game/search", map[string]interface{}{"NameIncludes":nameinclude, "Tag":tag, "SearchResults":res, "Error":err})
}

func HandleAccountPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w, r){
		user := models.RequestUser(r)
		games, err := models.Dal.GetGamesList(user.UserId)
		if err!=nil{
			panic(err)
		}
		user.Games = games
		friends, err := models.Dal.GetFriendsList(user.UserId)
		if err!=nil{
			panic(err)
		}
		user.Friends = friends
		models.RenderTemplate(w, r, "users/account", map[string]interface{}{"CurrentUser":user})
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
		friendarrstr := r.FormValue("accountFriendList")
		gamearrstr := r.FormValue("accountGameList")
		
		user, autherr := models.GlobalUserStore.Authenticate(username, oldPassword)
		if user==nil{
			user = models.RequestUser(r)
		}
		
		friends, err := models.Dal.GetFriendsList(user.UserId)
		if err!=nil{
			panic(err)
		}
		user.Friends = friends
		var friendarr []string
		err = json.Unmarshal([]byte(friendarrstr), &friendarr)
		if err!=nil || friendarr==nil{
			friendarr = []string{}
		}
		newfriendlist := []models.User{}
		delfriendlist := []string{}
		for _,friend := range user.Friends{
			found := false
			for _,id := range friendarr{
				if friend.UserId==id{
					found = true
				}
			}
			if found{
				newfriendlist = append(newfriendlist, friend)
			} else{
				delfriendlist = append(delfriendlist, friend.UserId)
			}
		}
		user.Friends = newfriendlist
		
		games, err := models.Dal.GetGamesList(user.UserId)
		if err!=nil{
			panic(err)
		}
		user.Games = games
		var gamearr []string
		err = json.Unmarshal([]byte(gamearrstr), &gamearr)
		if err!=nil || gamearr==nil{
			gamearr = []string{}
		}
		newgamelist := []models.Game{}
		delgamelist := []string{}
		for _,game := range user.Games{
			found := false
			for _,id := range gamearr{
				if game.GameId==id{
					found = true
				}
			}
			if found{
				newgamelist = append(newgamelist, game)
			} else{
				delgamelist = append(delgamelist, game.GameId)
			}
		}
		user.Games = newgamelist
		
		if autherr != nil {
			models.RenderTemplate(w, r, "users/account", map[string]interface{}{"Error": autherr.Error(), "CurrentUser":user})
			return
		}
		
		user.Email = email
		if newPassword != ""{
			if(newPassword != confirmNewPassword){
				models.RenderTemplate(w, r, "users/account", map[string]interface{}{"Error": "Passwords do not match.", "CurrentUser":user})
				return
			}
			user.SetPassword(newPassword)
		}
		
		err = models.Dal.UpdateUser(*user)
		if err!=nil{
			panic(err)
		}
		for _,remid := range delfriendlist{
			err = models.Dal.DeleteUserFriend(user.UserId, remid)
			if err!=nil{
				panic(err)
			}
		}
		for _,remid := range delgamelist{
			err = models.Dal.DeleteUserGame(user.UserId, remid)
			if err!=nil{
				panic(err)
			}
		}
		
		models.RenderTemplate(w, r, "users/account", map[string]interface{}{"CurrentUser":user})
	}
}

func HandleChatAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if models.SignedIn(w,r){
		//TODO: send the chatNewMsg to chatTo
		user := models.RequestUser(r)
		//toUser := r.FormValue("chatTo")
		chatID := r.FormValue("chatToID")
		theMessage := r.FormValue("chatNewMsg")
		chatto,err := models.Dal.FindUser(chatID)
		if err!=nil{
			panic(err)
		}
		models.Dal.SendMessage(*models.NewMessage(*user, *chatto, theMessage))

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
	game.UpdateRating()
	revs, err := models.Dal.FindTopReviewsByGame(gameid, 5)
	if err!=nil{
		panic(err)
	}
	vids, err := models.Dal.FindTopVideosByGame(gameid, 6)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r,"game/page", map[string]interface{}{"Game":game, "TopReviews":revs, "TopUserVids":vids})
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
func HandleGameClaimAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		err := models.Dal.AddUserGame(models.RequestUser(r).UserId, gameid)
		if err!=nil{
			panic(err)
		}
		http.Redirect(w,r, "/game/"+url.QueryEscape(gameid), http.StatusFound)
	}
}
func HandleReview(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	reviewid := params.ByName("wild")
	rev, err := models.Dal.FindReview(reviewid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r, "review/review", map[string]interface{}{"Review":rev})
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
		gameid := params.ByName("wild")
		//reviewid := r.FormValue("reviewID")
		user := models.RequestUser(r)
		//gameid := r.FormValue("gameID")
		review := r.FormValue("reviewBody")
		rating,_ := strconv.ParseFloat(r.FormValue("reviewRating"), 64)
		models.NewReview(user.UserId, gameid, review, rating)
	}
}
func HandleVideo(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	videoid := params.ByName("wild")
	vid, err := models.Dal.FindVideo(videoid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r, "video/uservideo", map[string]interface{}{"Video":vid, "VideoWidth":640, "VideoHeight":480})
}
func HandleUserVideos(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	userid := params.ByName("wild")
	user, err := models.GlobalUserStore.FindUser(userid)
	if err!=nil{
		panic(err)
	}
	var vids []models.Video
	vids, err = models.Dal.FindVideosByUser(user.UserId)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r, "videos/alluservideos", map[string]interface{}{"User":user, "AllVideos":vids})
}
func HandleGameVideos(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	gameid := params.ByName("wild")
	game, err := models.Dal.FindGame(gameid)
	if err!=nil{
		panic(err)
	}
	var vids []models.Video
	vids, err = models.Dal.FindVideosByGame(gameid)
	if err!=nil{
		panic(err)
	}
	models.RenderTemplate(w,r, "videos/allgamevideos", map[string]interface{}{"Game":game, "AllVideos":vids})
}

func HandleVideoNew(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		var gam *models.Game
		var err error
		gam, err = models.Dal.FindGame(gameid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "video/new", map[string]interface{}{"Game":gam})
	}
}
func HandleVideoNewAction(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	if models.SignedIn(w,r){
		gameid := params.ByName("wild")
		vidurl := r.FormValue("vidURL")
		user := models.RequestUser(r)
		vid := models.NewVideo(user.UserId, gameid, vidurl)
		err := models.Dal.CreateVideo(*vid)
		if err!=nil{
			panic(err)
		}
		http.Redirect(w,r, "/videos/"+url.QueryEscape(vid.ID), http.StatusFound)
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
	vids, err = models.Dal.FindTopVideosByUser(userid, 10)
	if err!=nil{
		panic(err)
	}
	revs, err = models.Dal.FindTopReviewsByUser(userid, 10)
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
		err = models.Dal.AddUserFriend(user.UserId, friendid)
		if err!=nil{
			panic(err)
		}
		models.RenderTemplate(w,r, "users/friend", map[string]interface{}{"Friend":friend})
	}
}
