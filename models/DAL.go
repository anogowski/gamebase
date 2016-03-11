package models

import (
	"database/sql"
	"errors"
	"html"
	"log"
	"os"
	"strconv"
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
	AddUserGame(userId, gameTitle string) error
	DeleteUserGame(userId, gameTitle string) error
	AddUserFriend(userId, friendId string) error
	DeleteUserFriend(userId, friendId string) error
	FindUser(id string) (*User, error)
	FindUserByName(name string) (*User, error)

	GetUsers() ([]User, error)
	SendMessage(message Message) error
	GetGamesList(userId string) ([]Game, error)
	GetFriendsList(userId string) ([]User, error)
	GetMessages(userId string) ([]Message, error)

	//GAME
	CreateGame(gam Game) error
	UpdateGame(game Game) error
	DeleteGame(gameId string) error
	FindGame(id string) (*Game, error)
	AddGameTag(gameId, tag string) error
	DeleteGameTag(gameId, tag string) error
	GetGames(amnt, skip int) ([]Game, error)
	SearchGames(search string) ([]Game, error)

	//Review
	CreateReview(review Review) error
	UpdateReview(review Review) error
	DeleteReview(reviewId string) error
	FindReview(reviewId string) (*Review, error)
	GetReviewsByGame(gameId string) ([]Review, error)
	GetReviewsGameCount(gameId string) (int, error)
	GetReviewsByUser(userId string) ([]Review, error)
	GetReviewsUserCount(userId string) (int, error)
	FindTopReviewsByGame(gameId string, amnt int) ([]Review, error)
	FindTopReviewsByUser(userId string, amnt int) ([]Review, error)

	//Tags
	CreateTag(tag string) error
	UpdateTag(oldTag, newTag string) error
	DeleteTag(tag string) error
	FindTag(tag string) error
	GetTags() ([]string, error)
	FindGamesByTag(tag string) ([]Game, error)
	FindTagsByGame(gameid string) ([]string, error)

	//Video
	CreateVideo(vid Video) error
	FindVideo(videoid string) (*Video, error)
	FindVideosByUser(userid string) ([]Video, error)
	FindVideosByGame(gameid string) ([]Video, error)
	FindTopVideosByUser(userid string, amnt int) ([]Video, error)
	FindTopVideosByGame(gameid string, amnt int) ([]Video, error)
	GetGameVideosCount(gameId string) (int, error)
	GetUserVideosCount(userId string) (int, error)
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
	if _, err = this.db.Exec("INSERT INTO users VALUES('" + user.UserId + "', '" + html.EscapeString(user.UserName) + "', '" + user.Password + "', '" + html.EscapeString(user.Email) + "')"); err != nil {
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
	user.UserName = html.UnescapeString(user.UserName)
	user.Email = html.UnescapeString(user.Email)
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
	user.UserName = html.UnescapeString(user.UserName)
	user.Email = html.UnescapeString(user.Email)
	return &user, nil
}

func (this *DataAccessLayer) UpdateUser(user User) error {
	if _, err := this.db.Exec("UPDATE users SET name='" + html.EscapeString(user.UserName) + "', password='" + user.Password + "', email='" + html.EscapeString(user.Email) + "' WHERE id='" + user.UserId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) AddUserGame(userId, gameId string) error {
	if _, err := this.db.Exec("INSERT INTO user_games VALUES('" + userId + "', '" + gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) DeleteUserGame(userId, gameId string) error {
	if _, err := this.db.Exec("DELETE FROM user_games WHERE (id='" + userId + "' AND gameId='" + gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) AddUserFriend(userId, friendId string) error {
	if _, err := this.db.Exec("INSERT INTO friends VALUES('" + userId + "', '" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) DeleteUserFriend(userId, friendId string) error {
	if _, err := this.db.Exec("DELETE FROM friends WHERE (id='" + userId + "' AND friendId='" + friendId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) SendMessage(message Message) error {
	user, err := this.FindUser(message.To.UserId)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("User does not exist")
	}
	if _, err = this.db.Exec("INSERT INTO messaging VALUES('" + message.From.UserId + "', '" + message.To.UserId + "', '" + html.EscapeString(message.TheMessage) + "', '" + time.Now().Format("2006-01-02 15:04:05") + "')"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) GetUsers() ([]User, error) {
	rows, err := this.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := []User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (this *DataAccessLayer) GetMessages(userId string) ([]Message, error) {

	rows, err := this.db.Query("SELECT * FROM messaging WHERE touserid ='" + userId + "' ORDER BY senttime")
	if err != nil {
		return nil, err
	}
	messages := []Message{}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.From, &message.To, &message.TheMessage, &message.TimeStamp)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (this *DataAccessLayer) CreateGame(gam Game) error {
	game, err := this.FindGame(gam.GameId)
	if err != nil {
		return err
	}
	if game != nil {
		return errors.New("Game already exists.")
	}
	if _, err = this.db.Exec("INSERT INTO games VALUES('" + gam.GameId + "', '" + html.EscapeString(gam.Title) + "', '" + html.EscapeString(gam.Developer) + "', '" + html.EscapeString(gam.Publisher) + "', '" + html.EscapeString(gam.Description) + "', '" + html.EscapeString(gam.URL) + "')"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) UpdateGame(game Game) error {
	if _, err := this.db.Exec("UPDATE games SET title='" + html.EscapeString(game.Title) + "', developer='" + html.EscapeString(game.Developer) + "', publisher='" + html.EscapeString(game.Publisher) + "', description='" + html.EscapeString(game.Description) + "', url='" + html.EscapeString(game.URL) + "' WHERE id='" + game.GameId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) DeleteGame(gameId string) error {
	if _, err := this.db.Exec("DELETE FROM games WHERE (id='" + gameId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) FindGame(id string) (*Game, error) {
	row := this.db.QueryRow("SELECT id,title,developer,publisher,description,url FROM games WHERE id='" + id + "'")
	game := Game{}
	err := row.Scan(&game.GameId, &game.Title, &game.Developer, &game.Publisher, &game.Description, &game.URL)
	game.Title = html.UnescapeString(game.Title)
	game.Developer = html.UnescapeString(game.Developer)
	game.Publisher = html.UnescapeString(game.Publisher)
	game.Description = html.UnescapeString(game.Description)
	game.URL = html.UnescapeString(game.URL)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &game, err
	}
	return &game, nil
}

func (this *DataAccessLayer) AddGameTag(gameId, tag string) error {
	if _, err := this.db.Exec("INSERT INTO game_tags VALUES('" + gameId + "', '" + html.EscapeString(tag) + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) DeleteGameTag(gameId, tag string) error {
	if _, err := this.db.Exec("DELETE FROM game_tags WHERE (gameId='" + gameId + "' AND tag='" + html.EscapeString(tag) + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) GetGames(amnt, skip int) ([]Game, error) {
	rows, err := this.db.Query("SELECT * FROM games LIMIT " + strconv.Itoa(amnt) + " OFFSET " + strconv.Itoa(skip))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	games := []Game{}
	defer rows.Close()
	for rows.Next() {
		var game Game
		err := rows.Scan(&game.GameId, &game.Title, &game.Developer, &game.Publisher, &game.Description, &game.URL)
		if err != nil {
			return games, err
		}
		game.Title = html.UnescapeString(game.Title)
		game.Developer = html.UnescapeString(game.Developer)
		game.Publisher = html.UnescapeString(game.Publisher)
		game.Description = html.UnescapeString(game.Description)
		game.URL = html.UnescapeString(game.URL)
		games = append(games, game)
	}
	return games, nil
}

func (this *DataAccessLayer) SearchGames(search string) ([]Game, error) {
	rows, err := this.db.Query("SELECT * FROM games WHERE lower(title) LIKE lower('%" + search + "%')")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	games := []Game{}
	defer rows.Close()
	for rows.Next() {
		var game Game
		err := rows.Scan(&game.GameId, &game.Title, &game.Developer, &game.Publisher, &game.Description, &game.URL)
		if err != nil {
			return games, err
		}
		game.Title = html.UnescapeString(game.Title)
		game.Developer = html.UnescapeString(game.Developer)
		game.Publisher = html.UnescapeString(game.Publisher)
		game.Description = html.UnescapeString(game.Description)
		game.URL = html.UnescapeString(game.URL)
		games = append(games, game)
	}
	return games, nil
}
func (this *DataAccessLayer) GetGamesList(userId string) ([]Game, error) {
	rows, err := this.db.Query("SELECT gameid,title,developer,publisher,description,url FROM user_games JOIN games ON user_games.gameid = games.id WHERE userid ='" + userId + "'")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	games := []Game{}
	defer rows.Close()
	for rows.Next() {
		var game Game
		err := rows.Scan(&game.GameId, &game.Title, &game.Developer, &game.Publisher, &game.Description, &game.URL)
		if err != nil {
			return games, err
		}
		game.Title = html.UnescapeString(game.Title)
		game.Developer = html.UnescapeString(game.Developer)
		game.Publisher = html.UnescapeString(game.Publisher)
		game.Description = html.UnescapeString(game.Description)
		game.URL = html.UnescapeString(game.URL)
		games = append(games, game)
	}
	return games, nil
}

func (this *DataAccessLayer) GetFriendsList(userId string) ([]User, error) {
	rows, err := this.db.Query("SELECT friendid,name,password,email FROM friends JOIN users ON friends.friendID = users.id WHERE friends.id ='" + userId + "'")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	users := []User{}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (this *DataAccessLayer) CreateReview(review Review) error {
	rview, err := this.FindReview(review.ReviewId)
	if err != nil {
		return err
	}
	if rview != nil {
		return errors.New("Review already exists.")
	}
	f := strconv.FormatFloat(review.Rating, 'g', 2, 64)
	if _, err = this.db.Exec("INSERT INTO reviews VALUES('" + review.ReviewId + "', '" + review.UserId + "', '" + review.GameId + "', '" + html.EscapeString(review.Body) + "', '" + f + "', " + strconv.Itoa(review.Likes) + ", " + strconv.Itoa(review.Dislikes) + ")"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) UpdateReview(review Review) error {
	f := strconv.FormatFloat(review.Rating, 'g', 2, 64)
	if _, err := this.db.Exec("UPDATE reviews SET body='" + html.EscapeString(review.Body) + "', rating='" + f + "', likes='" + strconv.Itoa(review.Likes) + "', dislikes='" + strconv.Itoa(review.Dislikes) + "' WHERE id='" + review.ReviewId + "'"); err != nil {
		return err
	}
	return nil
}

func (this *DataAccessLayer) DeleteReview(reviewId string) error {
	if _, err := this.db.Exec("DELETE FROM reviews WHERE (' reviewId=" + reviewId + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) FindReview(reviewId string) (*Review, error) {
	row := this.db.QueryRow("SELECT * FROM reviews WHERE reviewId='" + reviewId + "'")
	review := Review{}
	err := row.Scan(&review.ReviewId, &review.UserId, &review.GameId, &review.Body, &review.Rating, &review.Likes, &review.Dislikes)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &review, err
	}
	review.Body = html.UnescapeString(review.Body)
	return &review, nil
}

func (this *DataAccessLayer) GetReviewsByGame(gameId string) ([]Review, error) {
	rows, err := this.db.Query("SELECT * FROM reviews WHERE gameid ='" + gameId + "'")
	if err != nil {
		return nil, err
	}
	reivews := []Review{}
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.ReviewId, &review.UserId, &review.GameId, &review.Body, &review.Rating, &review.Likes, &review.Dislikes)
		if err != nil {
			return reivews, err
		}
		review.Body = html.UnescapeString(review.Body)
		reivews = append(reivews, review)
	}
	return reivews, nil
}

func (this *DataAccessLayer) GetReviewsByUser(userId string) ([]Review, error) {
	rows, err := this.db.Query("SELECT * FROM reviews WHERE userid ='" + userId + "'")
	if err != nil {
		return nil, err
	}
	reivews := []Review{}
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.ReviewId, &review.UserId, &review.GameId, &review.Body, &review.Rating, &review.Likes, &review.Dislikes)
		if err != nil {
			return reivews, err
		}
		review.Body = html.UnescapeString(review.Body)
		reivews = append(reivews, review)
	}
	return reivews, nil
}

func (this *DataAccessLayer) GetReviewsUserCount(userId string) (int, error) {
	rows, err := this.db.Query("SELECT COUNT(userid) FROM reviews WHERE userid ='" + userId + "'")
	if err != nil {
		return 0, err
	}
	numReviews := 0
	for rows.Next() {
		err = rows.Scan(&numReviews)
		if err != nil {
			return numReviews, err
		}
	}
	return numReviews, nil
}

func (this *DataAccessLayer) GetReviewsGameCount(gameId string) (int, error) {
	rows, err := this.db.Query("SELECT COUNT(gameid) FROM reviews WHERE gameid ='" + gameId + "'")
	if err != nil {
		return 0, err
	}
	numReviews := 0
	for rows.Next() {
		err = rows.Scan(&numReviews)
		if err != nil {
			return numReviews, err
		}
	}
	return numReviews, nil
}

func (this *DataAccessLayer) FindTopReviewsByGame(gameId string, amnt int) ([]Review, error) {
	revs := []Review{}
	rows, err := this.db.Query("SELECT * FROM reviews WHERE gameid ='" + gameId + "' ORDER BY likes DESC, dislikes ASC LIMIT " + strconv.Itoa(amnt))
	if err != nil {
		if err == sql.ErrNoRows {
			return revs, nil
		}
		return revs, err
	}
	defer rows.Close()
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.ReviewId, &review.UserId, &review.GameId, &review.Body, &review.Rating, &review.Likes, &review.Dislikes)
		if err != nil {
			return revs, err
		}
		review.Body = html.UnescapeString(review.Body)
		revs = append(revs, review)
	}
	return revs, nil
}
func (this *DataAccessLayer) FindTopReviewsByUser(userId string, amnt int) ([]Review, error) {
	revs := []Review{}
	rows, err := this.db.Query("SELECT * FROM reviews WHERE userid ='" + userId + "' ORDER BY likes DESC, dislikes ASC LIMIT " + strconv.Itoa(amnt))
	if err != nil {
		if err == sql.ErrNoRows {
			return revs, nil
		}
		return revs, err
	}
	defer rows.Close()
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.ReviewId, &review.UserId, &review.GameId, &review.Body, &review.Rating, &review.Likes, &review.Dislikes)
		if err != nil {
			return revs, err
		}
		review.Body = html.UnescapeString(review.Body)
		revs = append(revs, review)
	}
	return revs, nil
}

func (this *DataAccessLayer) CreateTag(tag string) error {
	err := this.FindTag(tag)
	if err == sql.ErrNoRows {
		if _, err = this.db.Exec("INSERT INTO tags VALUES('" + html.EscapeString(tag) + "')"); err != nil {
			return err
		}
	} else if err != nil {
		panic(err)
	}
	return nil
}

func (this *DataAccessLayer) DeleteTag(tag string) error {
	if _, err := this.db.Exec("DELETE FROM tags WHERE (' name=" + html.EscapeString(tag) + "')"); err != nil {
		return err
	}
	return nil

}

func (this *DataAccessLayer) FindTag(tag string) error {
	row := this.db.QueryRow("SELECT name FROM tags WHERE name='" + html.EscapeString(tag) + "'")
	err := row.Scan(&tag)
	return err
}

func (this *DataAccessLayer) UpdateTag(oldTag, newTag string) error {
	if _, err := this.db.Exec("UPDATE tags SET name='" + html.EscapeString(newTag) + "' WHERE name='" + html.EscapeString(oldTag) + "'"); err != nil {
		return err
	}
	return nil
}
func (this *DataAccessLayer) GetTags() ([]string, error) {
	rows, err := this.db.Query("SELECT name FROM tags")
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, html.UnescapeString(tag))
	}
	return tags, nil
}

func (this *DataAccessLayer) FindGamesByTag(tag string) ([]Game, error) {
	games := []Game{}
	rows, err := this.db.Query("SELECT gameid FROM game_tags WHERE tag='" + tag + "'")
	if err != nil {
		if err == sql.ErrNoRows {
			return games, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var gameid string
		err = rows.Scan(&gameid)
		if err != nil {
			return games, err
		}
		game, err := this.FindGame(gameid)
		if err != nil {
			return games, err
		}
		games = append(games, *game)
	}
	return games, nil
}
func (this *DataAccessLayer) FindTagsByGame(gameid string) ([]string, error) {
	tags := []string{}
	rows, err := this.db.Query("SELECT tag FROM game_tags WHERE gameid='" + gameid + "'")
	if err != nil {
		if err == sql.ErrNoRows {
			return tags, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, html.UnescapeString(tag))
	}
	return tags, nil
}

func (this *DataAccessLayer) CreateVideo(vid Video) error {
	if _, err := this.db.Exec("INSERT INTO videos VALUES('" + vid.ID + "', '" + vid.UserID + "', '" + vid.GameID + "', '" + html.EscapeString(vid.URL) + "', " + strconv.Itoa(vid.Likes) + ", " + strconv.Itoa(vid.Dislikes) + ")"); err != nil {
		return err
	}
	return nil
}
func (this *DataAccessLayer) FindVideo(videoid string) (*Video, error) {
	row := this.db.QueryRow("SELECT * FROM videos WHERE videoid='" + videoid + "'")
	vid := Video{}
	err := row.Scan(&vid.ID, &vid.UserID, &vid.GameID, &vid.URL, &vid.Likes, &vid.Dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	vid.URL = html.UnescapeString(vid.URL)
	return &vid, nil
}
func (this *DataAccessLayer) FindVideosByUser(userid string) ([]Video, error) {
	vids := []Video{}
	rows, err := this.db.Query("SELECT videoid,userid,gameid,url,likes,dislikes FROM videos WHERE userid='" + userid + "'")
	if err != nil {
		if err == sql.ErrNoRows {
			return vids, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vid := Video{}
		err := rows.Scan(&vid.ID, &vid.UserID, &vid.GameID, &vid.URL, &vid.Likes, &vid.Dislikes)
		if err != nil {
			return vids, err
		}
		vid.URL = html.UnescapeString(vid.URL)
		vids = append(vids, vid)
	}
	return vids, nil
}
func (this *DataAccessLayer) FindVideosByGame(gameid string) ([]Video, error) {
	vids := []Video{}
	rows, err := this.db.Query("SELECT videoid,userid,gameid,url,likes,dislikes FROM videos WHERE gameid='" + gameid + "'")
	if err != nil {
		if err == sql.ErrNoRows {
			return vids, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vid := Video{}
		err := rows.Scan(&vid.ID, &vid.UserID, &vid.GameID, &vid.URL, &vid.Likes, &vid.Dislikes)
		if err != nil {
			return vids, err
		}
		vid.URL = html.UnescapeString(vid.URL)
		vids = append(vids, vid)
	}
	return vids, nil
}
func (this *DataAccessLayer) FindTopVideosByUser(userid string, amnt int) ([]Video, error) {
	vids := []Video{}
	rows, err := this.db.Query("SELECT videoid,userid,gameid,url,likes,dislikes FROM videos WHERE userid='" + userid + "' ORDER BY likes DESC, dislikes ASC LIMIT " + strconv.Itoa(amnt))
	if err != nil {
		if err == sql.ErrNoRows {
			return vids, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vid := Video{}
		err := rows.Scan(&vid.ID, &vid.UserID, &vid.GameID, &vid.URL, &vid.Likes, &vid.Dislikes)
		if err != nil {
			return vids, err
		}
		vid.URL = html.UnescapeString(vid.URL)
		vids = append(vids, vid)
	}
	return vids, nil
}
func (this *DataAccessLayer) FindTopVideosByGame(gameid string, amnt int) ([]Video, error) {
	vids := []Video{}
	rows, err := this.db.Query("SELECT videoid,userid,gameid,url,likes,dislikes FROM videos WHERE gameid='" + gameid + "' ORDER BY likes DESC, dislikes ASC LIMIT " + strconv.Itoa(amnt))
	if err != nil {
		if err == sql.ErrNoRows {
			return vids, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vid := Video{}
		err := rows.Scan(&vid.ID, &vid.UserID, &vid.GameID, &vid.URL, &vid.Likes, &vid.Dislikes)
		if err != nil {
			return vids, err
		}
		vid.URL = html.UnescapeString(vid.URL)
		vids = append(vids, vid)
	}
	return vids, nil
}

func (this *DataAccessLayer) GetGameVideosCount(gameId string) (int, error) {
	rows, err := this.db.Query("SELECT COUNT(gameid) FROM videos WHERE gameid ='" + gameId + "'")
	if err != nil {
		return 0, err
	}
	numVideos := 0
	for rows.Next() {
		err = rows.Scan(&numVideos)
		if err != nil {
			return numVideos, err
		}
	}
	return numVideos, nil
}

func (this *DataAccessLayer) GetUserVideosCount(userId string) (int, error) {
	rows, err := this.db.Query("SELECT COUNT(userid) FROM videos WHERE userid ='" + userId + "'")
	if err != nil {
		return 0, err
	}
	numVideos := 0
	for rows.Next() {
		err = rows.Scan(&numVideos)
		if err != nil {
			return numVideos, err
		}
	}
	return numVideos, nil
}
