package models

type Game struct {
	Title     string
	Publisher string
	Rating    float64
	Review    []Review
}

const MAX_RATING float64 = 5

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

func (this *Game) UpdateRating(rating float64) {
	if rating > MAX_RATING {
		this.Rating = MAX_RATING
	} else if rating < 0 {
		this.Rating = 0
	} else {
		this.Rating = rating
	}
}

func (this *Game) UpdateReview(review Review) {

}
