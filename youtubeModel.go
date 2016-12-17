package youtime

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Content     string        `json:"content" bson:"content"`
	Time        int           `json:"time" bson:"time"`
	TimeCreated time.Time     `json:"timecreated" bson:"timecreated"`
}

type URL struct {
	Site string `json:"site" bson:"site"`
	ID   string `json:"id" bson:"id"`
}

type Video struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Url     URL           `json:"url" bson:"url"`
	Comment []Comment     `json:"comment" bson:"comment"`
}
