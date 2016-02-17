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
func NewSession(w http.ResponseWriter) *Session{
	exp := time.Now().Add(sessionLength)
	sess := &Session{
		ID : GenerateID("sess", sessionIDLength),
		Expiry : exp,
	}
	cookie := http.Cookie{
		Name: sessionCookieName,
		Value: sess.ID,
		Expires: exp,
	}
	http.SetCookie(w, &cookie)
	return sess
}
func RequestSession(r *http.Request) *Session{
	cookie, err := r.Cookie(sessionCookieName)
	if err!=nil{
		return nil
	}
	sess, err := globalSessionStore.Find(cookie.Value)
	if err!=nil{
		panic(err)
	}
	if sess==nil{
		return nil
	}
	if sess.Expired(){
		globalSessionStore.Delete(sess)
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
func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session{
	sess := RequestSession(r)
	if sess==nil{
		sess = NewSession(w)
	}
	return sess
}
func RequireLogin(w http.ResponseWriter, r *http.Request){
	if RequestUser(r)!=nil{
		return
	}
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}
