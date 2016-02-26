package models
import (
	"html/template"
	"net/http"
	"bytes"
	"html"
	"fmt"
)

var layoutFuncs = template.FuncMap{
	"yield":func()(string,error){
		return "",fmt.Errorf("bad yield called")
	},
	"yieldmenu":func()(string,error){
		return "",fmt.Errorf("bad yieldmenu called")
	},
	"yieldchat":func()(string,error){
		return "",fmt.Errorf("bad yieldchat called")
	},
	"RenderTemplateGameLink":func(href, onclick string)(string,error){
		return "",fmt.Errorf("bad RenderTemplateGameLink called")
	},
	"RenderTemplateRating":func(rating float64)(string,error){
		return "",fmt.Errorf("bad RenderTemplateRating called")
	},
	"RenderTemplateReview":func(rev Review, revclass string)(string,error){
		return "",fmt.Errorf("bad RenderTemplateReview called")
	},
	"RenderTemplateVideo":func(url, width, height interface{})(string,error){
		return "",fmt.Errorf("bad RenderTemplateVideo called")
	},
	"RenderTemplateUserVideo":func(vid Video, width, height interface{})(string,error){
		return "",fmt.Errorf("bad RenderTemplateUserVideo called")
	},
	"FindUserNameByID":func(userid string)(string,error){
		user, err := GlobalUserStore.FindUser(userid)
		if err!=nil{
			return "",err
		}
		return user.UserName,nil
	},
	"ReviewCount":func(userid string)(string,error){
		return "ReviewCount not yet implemented",nil
		//count, err := GlobalReviewStore.CountByUser(userid)
		//if err!=nil{
		//	return "0",err
		//}
		//return string(count),nil
	},
	"VideoCount":func(userid string)(string,error){
		return "VideoCount not yet implemented",nil
		//count, err := GlobalVideoStore.CountByUser(userid)
		//if err!=nil{
		//	return "0",err
		//}
		//return string(count),nil
	},
	"URLQueryEscaper":func(s interface{})(string,error){
		return template.URLQueryEscaper(s),nil
	},
}
var layout = template.Must(template.New("layout.html").Funcs(layoutFuncs).ParseFiles("templates/layout.html"))
var laytemplates = template.Must(template.New("t").Funcs(layoutFuncs).ParseFiles("templates/chat.html", "templates/menu.html"))
var templates = template.Must(template.New("t").Funcs(layoutFuncs).ParseGlob("templates/**/*.html"))
var errTemplate = `<h1>Error rendering template %s</h1><p>%s</p>`

func RenderTemplate(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}){
	if data==nil{
		data = map[string]interface{}{}
	}
	data["CurrentUser"] = RequestUser(r)
	data["Flash"] = r.URL.Query().Get("flash")
	data["Taglist"], _ = GlobalTagStore.GetTags()
	//data["ChatMessages"] = MessageStore.GetMessagesTo(data["CurrentUser"])
	
	funcs := template.FuncMap{
		"yield":func()(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, page, data)
			return template.HTML(buf.String()), err
		},
		"yieldmenu":func()(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := laytemplates.ExecuteTemplate(buf, "menu", data)
			return template.HTML(buf.String()), err
		},
		"yieldchat":func()(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := laytemplates.ExecuteTemplate(buf, "chat", data)
			return template.HTML(buf.String()), err
		},
		"RenderTemplateGameLink":func(href, onclick string)(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "home/index", map[string]interface{}{"GameLinkHREF":href, "GameLinkONCLICK":onclick})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateRating":func(rating float64)(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "review/rating", map[string]interface{}{"Rating":rating})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateReview":func(rev Review, revclass string)(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "review/review", map[string]interface{}{"ReviewClass":revclass, "Review":rev})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateVideo":func(url, width, height interface{})(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "game/video", map[string]interface{}{"VideoURL":url, "VideoWidth":width, "VideoHeight":height})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateUserVideo":func(vid Video, width, height interface{})(template.HTML,error){
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "users/video", map[string]interface{}{"Video":vid, "VideoWidth":width, "VideoHeight":height})
			return template.HTML(buf.String()), err
		},
		"FindUserNameByID":layoutFuncs["FindUserNameByID"],
		"ReviewCount":layoutFuncs["ReviewCount"],
		"VideoCount":layoutFuncs["VideoCount"],
		"URLQueryEscaper":layoutFuncs["URLQueryEscaper"],
	}
	layoutclone, _ := layout.Clone()
	layoutclone.Funcs(funcs)
	err := layoutclone.Execute(w, data)
	if err!=nil{
		http.Error(w, fmt.Sprintf(errTemplate, html.EscapeString(page), html.EscapeString(err.Error())), http.StatusInternalServerError)
	}
}
