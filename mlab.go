package youtime

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	URI        string
	Dbname     string
	Collection string
}

func CreateVideoMongo(item Video, mongo Mongodb) (Video, error) {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Video{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	thumbnail := "http://img.youtube.com/vi/" + item.Url.ID + "/0.jpg"
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	newVideo := Video{Id: bson.NewObjectId(), Url: item.Url, Comment: item.Comment, ThumbnailURL: thumbnail}
	err = collection.Insert(&newVideo)
	if err != nil {
		return Video{}, err
	}
	return newVideo, nil
}

func InsertCommentVideoMongo(id string, comment Comment, mongo Mongodb) (Video, error) {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Video{}, err
	}
	comment.TimeCreated = time.Now().UTC()
	comment.ID = bson.NewObjectId()
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)

	video := bson.ObjectIdHex(id)
	commentArray := bson.M{"$push": bson.M{"comment": bson.M{"$each": []Comment{comment}, "$sort": bson.M{"time": 1}}}}
	err = collection.UpdateId(video, commentArray)
	if err != nil {
		return Video{}, err
	}
	var result Video
	if bson.IsObjectIdHex(id) {
		err = collection.FindId(video).One(&result)
		if err != nil {
			return Video{}, err
		}
	} else {
		return Video{}, fmt.Errorf("Invalid input in ID %s", id)
	}
	return result, nil
}

func GetVideoByIdMongo(id string, mongo Mongodb) (Video, error) {
	var result Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Video{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)

	if bson.IsObjectIdHex(id) {
		err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
		if err != nil {
			return Video{}, err
		}
	} else {
		return Video{}, fmt.Errorf("Invalid input in ID %s", id)
	}

	return result, nil
}

func GetVideoByLinkMongo(url URL, mongo Mongodb) (Video, error) {
	var result Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Video{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	err = collection.Find(bson.M{"url": url}).One(&result)
	if err != nil {
		return Video{}, err
	}

	return result, nil
}

func GetAllVideoMongo(limit, offset string, mongo Mongodb) ([]Video, error) {
	var result []Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Video{}, err
	}
	var limitInt int
	var skipInt int
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Video{}, err
		}
	}

	if offset != "" {
		skipInt, err = strconv.Atoi(offset)
		if err != nil {
			return []Video{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	err = collection.Find(nil).Skip(skipInt).Limit(limitInt).All(&result)

	if err != nil {
		return []Video{}, err
	}

	return result, nil
}

func GetRandomVideoMongo(limit string, mongo Mongodb) ([]Video, error) {
	var result []Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Video{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	var limitInt int
	var skip int

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Video{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	count, err := collection.Count()
	if err != nil {
		return []Video{}, err
	}
	if count > limitInt {
		skip = rand.Intn(count - limitInt)
	} else {
		skip = rand.Intn(count)
	}

	err = collection.Find(nil).Skip(skip).Limit(limitInt).All(&result)

	if err != nil {
		return []Video{}, err
	}

	return result, nil
}
