package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	UserName string
	Password string
	UserId   string
	Email    string
	Games    []Game
	Messages []Message
	Friends  []User
}

const (
	hashcost  = 10
	userIDlen = 20
)

func NewUser(user_name, pass, mail string) *User {
	user := User{}
	user.InitUser(user_name, pass, mail)
	return &user
}
func (this *User) InitUser(user_name string, pass string, mail string) {
	this.UserName = user_name
	this.SetPassword(pass)
	this.UserId = GenerateID("user_", userIDlen)
	this.Email = mail
}

func (this *User) SetPassword(newPWord string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPWord), hashcost)
	this.Password = string(hash)
}

func (this *User) AddGame(game Game) {
	this.Games = append(this.Games, game)
	Dal.AddUserGame(this.UserId, game.GameId)
}

func (this *User) AddFriend(friend User) {
	temp, err := Dal.GetFriendsList(this.UserId)
	if err != nil {
		panic(err)
	}
	this.Friends = temp

}

func (this *User) GetMessages() {
	temp, err := Dal.GetMessages(this.UserId)
	if err != nil {
		panic(err)
	}
	this.Messages = temp
}

func (this *User) GetFriends() {
	temp, err := Dal.GetFriendsList(this.UserId)
	if err != nil {
		panic(err)
	}
	this.Friends = temp
}

func (this *User) GetGames() {
	temp, err := Dal.GetGamesList(this.UserId)
	if err != nil {
		panic(err)
	}
	this.Games = temp
}

func (this *User) CheckPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(pass)) == nil
}

func (this *User) UpdateUser(user_name, pass, mail string) {
	this.UserName = user_name
	this.Password = pass
	this.Email = mail
	Dal.UpdateUser(*this)
}

func (this *User) AddGameToList(gameId string) {
	Dal.AddUserGame(this.UserId, gameId)
}

func (this *User) AddFriendToList(friendId string) {
	Dal.AddUserFriend(this.UserId, friendId)
}

func (this *User) DeleteGameFromList(gameId string) {
	Dal.DeleteUserGame(this.UserId, gameId)
}

func (this *User) DeleteFriendFromList(friendId string) {
	Dal.DeleteUserFriend(this.UserId, friendId)
}

func (this *User) SendMessage(message Message) {
	Dal.SendMessage(message)
}
