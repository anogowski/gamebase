package models

type Review struct {
	ReviewId string
	UserId   string
	GameId   string
	Body     string
	Rating   float64
	Likes    int
	Dislikes int
}

const REVIEW_MAX_RATING float64 = 5
const REVIEW_ID_LEN = 20

func NewReview(userId, gameId, body string, rating float64)*Review{
	rev := &Review{}
	rev.InitReview(userId, gameId, body, rating)
	return rev
}

func (this *Review) InitReview(userId, gameId, body string, rating float64) {
	this.Body = body
	this.UserId = userId
	this.GameId = gameId
	this.Rating = rating
	this.ReviewId = GenerateID("review_", REVIEW_ID_LEN)
	Dal.CreateReview(*this)
}

func (this *Review) UpdateReview(body, url string, rating float64) {
	this.Body = body
	if rating > REVIEW_MAX_RATING {
		this.Rating = REVIEW_MAX_RATING
	} else if rating < 0 {
		this.Rating = 0
	} else {
		this.Rating = rating
	}
	Dal.UpdateReview(*this)
}

func (this *Review) DeleteReview() {
	Dal.DeleteReview(this.ReviewId)
}
