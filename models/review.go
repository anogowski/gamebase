package models

type Review struct {
	ReviewId string
	UserId   string
	GameId   string
	Body     string
	Rating   float64
}

const REVIEW_MAX_RATING float64 = 5
const REVIEW_ID_LEN = 20

func (this *Review) InitReview(userId, gameId, body string, rating float64) {
	this.Body = body
	this.UserId = userId
	this.GameId = gameId
	this.Rating = rating
	this.ReviewId = GenerateID("review_", REVIEW_ID_LEN)
	err = Dal.CreateReview(*this)
}

func (this *Review) UpdateReview(body, url string, rating float64) {
	this.Body = body
	this.URL = url
	if rating > REVIEW_MAX_RATING {
		this.Rating = REVIEW_MAX_RATING
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
