package models

type Review struct {
	Title  string
	Body   string
	Rating float64
	UserId string
	GameId string
	URL    string
}

const MAX_RATING float64 = 5

func (this *Review) InitReview(title, body, userId string, rating float64, url string) {
	this.Title = title
	this.Body = body
	this.UserId = userId
	this.Rating = rating
	this.URL = url
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
