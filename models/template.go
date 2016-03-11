package models

import (
	"html/template"
	"encoding/json"
	"net/http"
	"strconv"
	"bytes"
	"html"
	"fmt"
)

func StringIndexOf(str string, needle string, start int) int {
	length := len(needle)
	max := len(str) - length
	for at := start; at < max; at++ {
		if needle == str[at:at+length] {
			return at
		}
	}
	return -1
}

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("bad yield called")
	},
	"yieldmenu": func() (string, error) {
		return "", fmt.Errorf("bad yieldmenu called")
	},
	"yieldchat": func() (string, error) {
		return "", fmt.Errorf("bad yieldchat called")
	},
	"RenderTemplateGameLink": func(game Game, href, onclick string) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateGameLink called")
	},
	"RenderTemplateRating": func(rating float64) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateRating called")
	},
	"RenderTemplateReview": func(rev Review) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateReview called")
	},
	"RenderTemplateVideo": func(url, width, height interface{}) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateVideo called")
	},
	"RenderTemplateGameVideo": func(url, width, height interface{}) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateGameVideo called")
	},
	"RenderTemplateUserVideo": func(vid Video, width, height interface{}) (string, error) {
		return "", fmt.Errorf("bad RenderTemplateUserVideo called")
	},
	"FindUserNameByID": func(userid string) (string, error) {
		user, err := GlobalUserStore.FindUser(userid)
		if err != nil {
			return "", err
		}
		return user.UserName, nil
	},
	"FindGameNameByID": func(gameid string) (string, error) {
		game, err := Dal.FindGame(gameid)
		if err != nil {
			return "", err
		}
		return game.Title, nil
	},
	"ReviewCount": func(userid string) (string, error) {
		count, err := Dal.GetReviewsUserCount(userid)
		if err != nil {
			return "0", err
		}
		return strconv.Itoa(count), nil
	},
	"VideoCount": func(userid string) (string, error) {
		count, err := Dal.GetUserVideosCount(userid)
		if err != nil {
			return "0", err
		}
		return strconv.Itoa(count), nil
	},
	"ReviewGameCount": func(gameid string) (string, error) {
		count, err := Dal.GetReviewsGameCount(gameid)
		if err != nil {
			return "0", err
		}
		return strconv.Itoa(count), nil
	},
	"VideoGameCount": func(gameid string) (string, error) {
		count, err := Dal.GetGameVideosCount(gameid)
		if err != nil {
			return "0", err
		}
		return strconv.Itoa(count), nil
	},
	"URLQueryEscaper": func(s interface{}) (string, error) {
		return template.URLQueryEscaper(s), nil
	},
	"JSONify": func(s interface{}) (string, error) {
		js, err := json.Marshal(s)
		return string(js), err
	},
	"HTMLnewlines": func(s string) (template.HTML, error) {
		buf := bytes.NewBuffer(nil)
		lastind := 0
		for ind := StringIndexOf(s, "\n", lastind); ind >= 0; ind = StringIndexOf(s, "\n", ind+1) {
			fmt.Fprintln(buf, template.HTMLEscapeString(s[lastind:ind])+"<br/>")
			lastind = ind
		}
		return template.HTML(buf.String()), nil
	},
}
var layout = template.Must(template.New("layout.html").Funcs(layoutFuncs).ParseFiles("templates/layout.html"))
var laytemplates = template.Must(template.New("t").Funcs(layoutFuncs).ParseFiles("templates/chat.html", "templates/menu.html"))
var templates = template.Must(template.New("t").Funcs(layoutFuncs).ParseGlob("templates/**/*.html"))
var errTemplate = `<h1>Error rendering template %s</h1><p>%s</p>`

func RenderTemplate(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	if _, ok := data["CurrentUser"]; !ok || (data["CurrentUser"].(*User) == nil) {
		data["CurrentUser"] = RequestUser(r)
	}
	data["Flash"] = r.URL.Query().Get("flash")
	data["Taglist"], _ = Dal.GetTags()
	if data["CurrentUser"].(*User) != nil {
		data["friendsList"], _ = Dal.GetFriendsList(data["CurrentUser"].(*User).UserId)
		data["ChatMessages"], _ = Dal.GetMessages(data["CurrentUser"].(*User).UserId)
	}
	var templateClone *template.Template

	renderFuncs := template.FuncMap{
		"RenderTemplateGameLink": func(game Game, href, onclick string) (template.HTML, error) {
			if href == "" {
				href = "/game/" + game.GameId
			}
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "game/gamelink", map[string]interface{}{"GameLinkHREF": href, "GameLinkONCLICK": onclick, "Game": game})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateRating": func(rating float64) (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "review/rating", map[string]interface{}{"Rating": rating})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateReview": func(rev Review) (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "review/review", map[string]interface{}{"Review": rev})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateVideo": func(url, width, height interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "video/embed", map[string]interface{}{"VideoURL": url, "VideoWidth": width, "VideoHeight": height})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateGameVideo": func(vid Video, width, height interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "video/gamevideo", map[string]interface{}{"Video": vid, "VideoWidth": width, "VideoHeight": height})
			return template.HTML(buf.String()), err
		},
		"RenderTemplateUserVideo": func(vid Video, width, height interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, "video/uservideo", map[string]interface{}{"Video": vid, "VideoWidth": width, "VideoHeight": height})
			return template.HTML(buf.String()), err
		},
	}
	templateClone, _ = templates.Clone()
	templateClone.Funcs(renderFuncs)

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templateClone.ExecuteTemplate(buf, page, data)
			return template.HTML(buf.String()), err
		},
		"yieldmenu": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := laytemplates.ExecuteTemplate(buf, "menu", data)
			return template.HTML(buf.String()), err
		},
		"yieldchat": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := laytemplates.ExecuteTemplate(buf, "chat", data)
			return template.HTML(buf.String()), err
		},
		"FindUserNameByID": layoutFuncs["FindUserNameByID"],
		"ReviewCount":      layoutFuncs["ReviewCount"],
		"VideoCount":       layoutFuncs["VideoCount"],
		"URLQueryEscaper":  layoutFuncs["URLQueryEscaper"],
		"JSONify":          layoutFuncs["JSONify"],
		"HTMLnewlines":     layoutFuncs["HTMLnewlines"],
	}
	layoutclone, _ := layout.Clone()
	layoutclone.Funcs(funcs).Funcs(renderFuncs)
	err := layoutclone.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf(errTemplate, html.EscapeString(page), html.EscapeString(err.Error())), http.StatusInternalServerError)
	}
}
