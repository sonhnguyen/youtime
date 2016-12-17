package youtime

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	Content     string    `json:"content" bson:"content"`
	Time        int       `json:"time" bson:"time"`
	DateCreated time.Time `json:"datecreated" bson:"timeupdated"`
}

type URL struct {
	Site string `json:"site" bson:"site"`
	Link string `json:"link" bson:"link"`
}

type Video struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Url     URL           `json:"url" bson:"url"`
	Comment []Comment     `json:"comment" bson:"comment"`
}
