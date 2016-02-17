package main

import (
	"review"
)

type Game struct {
	var title     string
	var publisher string
	var rating    float64
	var review    Review
	var url       

	func InitGame(tit string) {
		title = tit
	}

	func InitGame(tit string, string pub) {
		title = tit
		publisher = pub
	}

}
