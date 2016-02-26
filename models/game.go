package models

type Game struct {
	GameId    string
	Title     string
	Publisher string
	Rating    float64
	Review    []Review
}

const (
	MAX_RATING  float64 = 5
	GAME_ID_LEN         = 20
)

func NewGame(gameId, title, publisher string) *Game {
	game := Game{}
	game.InitGame(title, publisher)
	return &game
}

func (this *Game) InitGame(title string, pub string) {
	this.Title = title
	this.Publisher = pub
	this.GameId = GenerateID("game_", GAME_ID_LEN)
}

func (this *Game) UpdateTitle(title string) {
	this.Title = title
}

func (this *Game) UpdatePublisher(pub string) {
	this.Publisher = pub
}

func (this *Game) UpdateRating(rating float64) {
	//Get raitings from reviews
}
