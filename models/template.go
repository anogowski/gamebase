package models
import (
	"html/template"
	"net/http"
	"bytes"
	"html"
	"fmt"
)

var layoutFuncs = template.FuncMap{
	"URLQueryEscaper":func(s interface{})(string,error){
		return template.URLQueryEscaper(s),nil
	},
	"yield":func()(string,error){
		return "",fmt.Errorf("bad yield called")
	},
	"yieldmenu":func()(string,error){
		return "",fmt.Errorf("bad yieldmenu called")
	},
	"yieldchat":func()(string,error){
		return "",fmt.Errorf("bad yieldchat called")
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
	
	funcs := template.FuncMap{
		"URLQueryEscaper":func(s interface{})(string,error){
			return template.URLQueryEscaper(s),nil
		},
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
	}
	layoutclone, _ := layout.Clone()
	layoutclone.Funcs(funcs)
	err := layoutclone.Execute(w, data)
	if err!=nil{
		http.Error(w, fmt.Sprintf(errTemplate, html.EscapeString(page), html.EscapeString(err.Error())), http.StatusInternalServerError)
	}
}
