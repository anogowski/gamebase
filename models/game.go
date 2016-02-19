package models

type Game struct {
	Title     string
	Publisher string
	Rating    float64
	Review    []Review
	GameId    string
}

const MAX_RATING int = 5

func (this *Game) InitGame(title string, pub string) {
	this.Title = title
	this.Publisher = pub
}

func (this *Game) UpdateTitle(title string) {
	this.Title = title
}

func (this *Game) UpdatePublisher(pub string) {
	this.Publisher = pub
}

func (this *Game) UpdateRating(rating int) {
	if rating > MAX_RATING {
		this.Rating = MAX_RATING
	} else if raiting < 0 {
		this.Rating = 0
	} else {
		this.Rating = raiting
	}
}

func (this *Game) UpdateReview(review Review) {

}
