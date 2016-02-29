package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "gamebase/Godeps/_workspace/src/github.com/lib/pq"
)

var Dal DAL

func init() {
	Dal = NewDataAccessLayer()
}

type DataAccessLayer struct {
	db *sql.DB
}

func NewDataAccessLayer() *DataAccessLayer {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening the database: %q", err)
	}
	dal := &DataAccessLayer{db: db}
	return dal
}

type DAL interface {
	//USER
	CreateUser(name, pass, email string) (*User, error)
	UpdateUser(user User) error
	AddUserGame(user User, gameTitle string) error
	DeleteUserGame(user User, gameTitle string) error
	AddUserFriend(user User, friendId string) error
	DeleteUserFriend(user User, friendId string) error
	FindUser(id string) (*User, error)
	FindUserByName(name string) (*User, error)

	//GetUsers()
	SendMessage(from, to, message string) error
	//GetGamesList()
	//GetFriendsList()
	//GetMessages()

	//GAME
	CreateGame(id, title, publisher, url string) (*Game, error)
	UpdateGame(game Game) error
	DeleteGame(gameId string) error
	FindGame(id string) (*Game, error)

	//GetGames()

	//Review
	/*
		CreateReview(title, body, url, userId, gameId string, rating float64) (*Review, error)
		UpdateReview(title, body, url, userId, gameId string, rating float64) error
		DeleteReview(userId, gameId string) error
		FindReview(userId, gameId) (*Review, error)
		GetReviews()
	*/
	//Tags

	CreateTag(tag string) error
	UpdateTag(oldTag, newTag string) error
	//	RemoveTag()
	FindTag(tag string) error
	//	GetTags()

}

func (this *DataAccessLayer) CreateUser(name, pass, email string) (*User, error) {
	user, err := this.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Username already taken.")
	}
	user = NewUser(name, pass, email)
	if _, err = this.db.Exec("INSERT INTO users VALUES('" + user.UserId + "', '" + user.UserName + "', '" + user.Password + "', '" + user.Email + "')"); err != nil {
		return user, err
	}
	return user, nil
}
func (this *DataAccessLayer) FindUser(id string) (*User, error) {
	row := this.db.QueryRow("SELECT * FROM users WHERE id='" + id + "'")
	user := User{}
	err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &user, err
	}
	return &user, nil
}
func (this *DataAccessLayer) FindUserByName(name string) (*User, error) {
	row := this.db.QueryRow("SELECT id, name, password, email FROM users WHERE name='" + name + "'")
	user := User{}
	err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return &user, nil
}

func (this *DataAccessLayer) UpdateUser(user User) error {
	if _, err := this.db.Exec("UPDATE users SET name='" + user.UserName + "', password='" + user.Password + "', email='" + user.Email + "' WHERE id='" + user.UserId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) AddUserGame(user User, gameId string) error {
	if _, err := this.db.Exec("INSERT INTO user_games VALUES('" + user.UserId + "', '" + gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) DeleteUserGame(user User, gameId string) error {
	if _, err := this.db.Exec("DELETE FROM user_games WHERE (' id=" + user.UserId + "'AND gameId='" + gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) AddUserFriend(user User, friendId string) error {
	if _, err := this.db.Exec("INSERT INTO friends VALUES('" + user.UserId + "', '" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) DeleteUserFriend(user User, friendId string) error {
	if _, err := this.db.Exec("DELETE FROM friends WHERE (' id=" + user.UserId + "'AND friendId='" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) SendMessage(from, to, message string) error {
	user, err := FindUser(to)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("User does not exist")
	}
	if _, err = this.db.Exec("INSERT INTO messaging VALUES('" + from + "', '" + to + "', '" + message + "', '" + time.Time.Now() + "')"); err != nil {
		return game, err
	}
	return game, nil
}

func (this *DataAccessLayer) CreateGame(id, title, publisher, url string) (*Game, error) {
	game, err := this.FindGame(id)
	if err != nil {
		return nil, err
	}
	if game != nil {
		return nil, errors.New("Game already exists.")
	}
	game = NewGame(title, publisher, url)
	if _, err = this.db.Exec("INSERT INTO games VALUES('" + game.GameId + "', '" + game.Title + "', '" + game.Publisher + "', '" + game.URL + "')"); err != nil {
		return game, err
	}
	return game, nil
}

func (this *DataAccessLayer) UpdateGame(game Game) error {
	if _, err := this.db.Exec("UPDATE games SET title='" + game.Title + "', publisher='" + game.Publisher + "', url='" + game.URL + "' WHERE id='" + game.GameId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) DeleteGame(gameId string) error {
	if _, err := this.db.Exec("DELETE FROM games WHERE (' id=" + user.gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) FindGame(id string) (*Game, error) {
	row := this.db.QueryRow("SELECT * FROM games WHERE id='" + id + "'")
	game := Game{}
	err := row.Scan(&game.GameId, &game.Title, &game.Publisher, &game.Rating)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &game, err
	}
	return &game, nil
}

func (this *DataAccessLayer) CreateTag(tag string) error {
	err := this.FindTag(tag)
	if err == sql.ErrNoRows {
		if _, err = this.db.Exec("INSERT INTO labels VALUES('" + tag + "')"); err != nil {
			return err
		}
	} else if err != nil {
		panic(err)
	}
	return nil
}

func (this *DataAccessLayer) FindTag(tag string) error {
	row := this.db.QueryRow("SELECT name FROM labels WHERE name='" + tag + "'")
	err := row.Scan(&tag)
	return err
}
