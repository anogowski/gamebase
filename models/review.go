package models

type Review struct {
	Title  string
	Body   string
	URL    string
	UserId string
	GameId string
	Rating float64
}

const MAX_RATING float64 = 5

func (this *Review) InitReview(title, body, url, userId, gameId string, rating float64) {
	this.Title = title
	this.Body = body
	this.URL = url
	this.UserId = userId
	this.GameId = gameId
	this.Rating = rating

	//CALL DAL
}

func (this *Review) UpdateReview(title, body, url string, rating float64) {
	this.Title = title
	this.Body = body
	this.URL = url
	if rating > MAX_RATING {
		this.Rating = MAX_RATING
	} else if rating < 0 {
		this.Rating = 0
	} else {
		this.Rating = rating
	}
	//CALL DAL
}

func (this *Review) DeleteReview() {
	//CALL DAL
}
