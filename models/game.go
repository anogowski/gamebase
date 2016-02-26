package models

type Game struct {
	GameId    string
	Title     string
	Publisher string
	Rating    float64
	URL       string
	Review    []Review
}

const (
	GAME_MAX_RATING float64 = 5
	GAME_ID_LEN             = 20
)

func NewGame(gameId, title, publisher string) *Game {
	game := Game{}
	game.InitGame(title, publisher)
	return &game
}

func (this *Game) InitGame(title, pub, url string) {
	this.Title = title
	this.Publisher = pub
	this.URL = url
	this.GameId = GenerateID("game_", GAME_ID_LEN)
	Dal.CreateGame(this.GameId, this.Title, this.Publisher, this.URL)
}

func (this *Game) UpdateGame(title, pub, url string) {
	this.Title = title
	this.Publisher = pub
	this.URL = url
	Dal.UpdateGame(this)
}

func (this *Game) UpdateRating(rating float64) {
	//Get raitings from reviews
}
