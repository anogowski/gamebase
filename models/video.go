package models
import (
	
)

type Video struct{
	ID string
	UserID string
	GameID string
	URL string
	Likes int
	Dislikes int
}

const(
	videoidlen = 20
)

func NewVideo(userid, gameid, url string) *Video{
	id := GenerateID("vid_", videoidlen)
	return &Video{ID:id, UserID:userid, GameID:gameid, URL:url, Likes:0, Dislikes:0}
}
