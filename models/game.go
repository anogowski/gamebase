package main

import (
	
)

type Game struct {
	title     string
	publisher string
	rating    float64
	review    Review
	url       string
}

func (this *Game) InitGame(tit string, pub string) {
	this.title = tit
	this.publisher = pub
}
