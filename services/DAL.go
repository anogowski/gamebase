//TODO:	placeholder methods connecting to the database
package services

import (
	"database/sql"
	"errors"

	_ "gamebase/Godeps/_workspace/src/github.com/lib/pq"
)

type DAL interface {
	//USER
	CreateUser(name, pass, email string) (*User, error)
	FindUser(id string) (*User, error)
	FindUserByName(name string) (*User, error)
	UpdateUser(user User) error
	AddUserGame(user User, gameTitle string) error
	DeleteUserGame(user User, gameTitle string) error
	AddUserFriend(user User, friendId string) error
	DeleteUserFriend(user User, friendId string) error
	//GAME
	CreateGame(title, publisher string) (*Game, error)
	FindGame(id string) (*Game, error)
	UpdateGame(title, publisher string) error
	//Review
	CreateReview(name, pass, email string, rating float64) (*Review, error)
	UpdateReview()
	DeleteReview(userId, gameId string) error
}

func (this *DAL) CreateUser(name, pass, email string) (*User, error) {
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
func (this *DAL) FindUser(id string) (*User, error) {
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
func (this *DAL) FindUserByName(name string) (*User, error) {
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

func (this *DAL) UpdateUser(user User) error {
	if _, err := this.db.Exec("UPDATE users SET name='" + user.UserName + "', password='" + user.Password + "', email='" + user.Email + "' WHERE id='" + user.UserId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DAL) AddUserGame(user User, gameTitle string) error {
	if _, err := this.db.Exec("INSERT INTO user_games VALUES('" + user.UserId + "', '" + gameTitle + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DAL) DeleteUserGame(user User, gameTitle string) error {
	if _, err := this.db.Exec("DELETE FROM user_games WHERE (' id=" + user.UserId + "'AND gameTitle='" + gameTitle + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DAL) AddUserFriend(user User, friendId string) error {
	if _, err := this.db.Exec("INSERT INTO user_friends VALUES('" + user.UserId + "', '" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DAL) DeleteUserFriend(user User, friendId string) error {
	if _, err := this.db.Exec("DELETE FROM user_friends WHERE (' id=" + user.UserId + "'AND friendId='" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DAL) CreateGame(id, title, publisher string, rating float64) (*Game, error) {
	game, err := this.FindGame(id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Game already exists.")
	}
	game = NewGame(title, publisher)
	if _, err = this.db.Exec("INSERT INTO games VALUES('" + game.GameId + "', '" + game.Title + "', '" + game.Publisher + "', '" + game.Raiting + "')"); err != nil {
		return game, err
	}
	return game, nil
}

func (this *DAL) FindGame(id string) (*Game, error) {
	row := this.db.QueryRow("SELECT * FROM games WHERE id='" + id + "'")
	game := Game{}
	err := row.Scan(&game.GameId, &game.Title, &game.Publisher, &game.Raiting)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &game, err
	}
	return &game, nil
}
