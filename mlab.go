package youtime

import (
	"fmt"

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
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	newVideo := Video{Id: bson.NewObjectId(), Url: item.Url}
	collection.Insert(&newVideo)
	return newVideo, nil
}

func InsertCommentVideoMongo(id string, comment Comment, mongo Mongodb) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)

	video := bson.ObjectIdHex(id)
	commentArray := bson.M{"$push": bson.M{"comment": bson.M{"$each": []Comment{comment}, "$sort": bson.M{"time": 1}}}}
	err = collection.UpdateId(video, commentArray)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoByIdMongo(id string, mongo Mongodb) (Video, error) {
	var result Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		return Video{}, nil
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, nil
}

func GetVideoByLinkMongo(url URL, mongo Mongodb) (Video, error) {
	var result Video
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		return Video{}, nil
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	collection.Find(bson.M{"url": url}).One(&result)
	return result, nil
}
