package models

type Game struct {
	GameId      string
	Title       string
	Publisher   string
	Developer   string
	Rating      float64
	Description string
	URL         string
	Review      []Review
}

const (
	GAME_MAX_RATING float64 = 5
	GAME_ID_LEN             = 20
)

func NewGame(title, publisher, url string) *Game {
	game := Game{}
	game.InitGame(title, publisher, url)
	return &game
}

func (this *Game) InitGame(title, pub, url string) {
	this.Title = title
	this.Publisher = pub
	this.URL = url
	this.GameId = GenerateID("game_", GAME_ID_LEN)
	//Dal.CreateGame(this.GameId, this.Title, this.Publisher, this.URL)
}

func (this *Game) UpdateGame(title, pub, url string) {
	this.Title = title
	this.Publisher = pub
	this.URL = url
	Dal.UpdateGame(*this)
}

func (this *Game) DeleteGame() {
	Dal.DeleteGame(this.GameId)
}

func (this *Game) GetReviews() {
	temp, err := Dal.GetReviewsByGame(this.GameId)
	if err != nil {
		panic(err)
	}
	this.Review = temp
}

func (this *Game) UpdateRating() {

	reviews, err := Dal.GetReviewsByGame(this.GameId)
	if err != nil {
		panic(err)
	}

	var numReviews int = len(reviews)
	if numReviews > 0 {

		var sumRating float64 = 0.0
		for i := 0; i < numReviews; i++ {
			sumRating += reviews[i].Rating
		}
		this.Rating = sumRating / float64(numReviews)
	} else {
		this.Rating = 0.0
	}

}
