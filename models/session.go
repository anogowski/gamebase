package models
import (
	"time"
	"net/url"
	"net/http"
)

const (
	sessionLength = 24 * time.Hour
	sessionCookieName = "GamebaseSession"
	sessionIDLength = 20
)

type Session struct{
	ID	string
	UserID string
	Expiry time.Time
}
func (this *Session) Expired()bool{
	return this.Expiry.Before(time.Now())
}
func NewSession(w http.ResponseWriter, userid string) *Session{
	exp := time.Now().Add(sessionLength)
	sess := &Session{
		ID : GenerateID("sess_", sessionIDLength),
		UserID : userid,
		Expiry : exp,
	}
	cookie := http.Cookie{
		Name: sessionCookieName,
		Value: sess.ID,
		Expires: exp,
	}
	err := GlobalSessionStore.Save(sess)
	if err!=nil{
		panic(err)
	}
	http.SetCookie(w, &cookie)
	return sess
}
func RequestSession(r *http.Request) *Session{
	cookie, err := r.Cookie(sessionCookieName)
	if err!=nil{
		return nil
	}
	sess, err := GlobalSessionStore.Find(cookie.Value)
	if err!=nil{
		panic(err)
	}
	if sess==nil{
		return nil
	}
	if sess.Expired(){
		GlobalSessionStore.Delete(sess)
		return nil
	}
	return sess
}
func RequestUser(r *http.Request) *User{
	sess := RequestSession(r)
	if sess==nil || sess.UserID==""{
		return nil
	}
	user, err := GlobalUserStore.FindUser(sess.UserID)
	if err!=nil{
		panic(err)
	}
	return user
}
func FindOrCreateSession(w http.ResponseWriter, r *http.Request, userid string) *Session{
	sess := RequestSession(r)
	if sess==nil{
		sess = NewSession(w, userid)
	}
	return sess
}
func SignedIn(w http.ResponseWriter, r *http.Request)bool{
	if RequestUser(r)!=nil{
		return true
	}
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
	return false
}
