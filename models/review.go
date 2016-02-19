package models

type Review struct {
	Title  string
	Body   string
	Rating float64
	UserId string
	URL    string
}

func (this *Review) InitReview(title, body, userId string, rating float64) {
	this.Title = title
	this.Body = body
	this.UserId = userId
	this.Rating = rating
}

func (this *Review) InitReview(title, body, userId string, rating float64, url string) {
	this.Title = title
	this.Body = body
	this.UserId = userId
	this.Rating = rating
	this.URL = url
}

func (this *Review) UpdateTitle(title string) {
	this.Title = title
}

func (this *Review) UpdateBody(body string) {
	this.Body = body
}

func (this *Review) UpdateURL(url string) {
	this.URL = url
}
