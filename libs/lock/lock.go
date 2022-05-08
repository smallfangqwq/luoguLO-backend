package lock

import (
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

type LockDBTemplate struct {
	PostID     int
	LockStatus int
}

func StartLock(session *mgo.Session, postID int) {
	var collection = session.DB("luogulo").C("discussLock")
	collectionCount, err := collection.Find(bson.M{"postid": postID}).Count()
	if err != nil {
		panic(err)
	}
	// new
	var saveVersion LockDBTemplate
	saveVersion.PostID = postID
	saveVersion.LockStatus = 1
	if collectionCount == 0 {
		collection.Insert(saveVersion)
	} else {
		collection.Update(bson.M{"postid": postID}, saveVersion)
	}
}

func CloseLock(session *mgo.Session, postID int) {
	var collection = session.DB("luogulo").C("discussLock")
	collectionCount, err := collection.Find(bson.M{"postid": postID}).Count()
	if err != nil {
		panic(err)
	}
	var saveVersion LockDBTemplate
	saveVersion.PostID = postID
	saveVersion.LockStatus = 0
	if collectionCount == 0 {
		collection.Insert(saveVersion)
	} else {
		collection.Update(bson.M{"postid": postID}, saveVersion)
	}
}

func SetLockStatus(session *mgo.Session, postID int, status int) {
	var collection = session.DB("luogulo").C("discussLock")
	collectionCount, err := collection.Find(bson.M{"postid": postID}).Count()
	if err != nil {
		panic(err)
	}
	var saveVersion LockDBTemplate
	saveVersion.PostID = postID
	saveVersion.LockStatus = status
	if collectionCount == 0 {
		collection.Insert(saveVersion)
	} else {
		collection.Update(bson.M{"postid": postID}, saveVersion)
	}
}

func GetLockStatus(session *mgo.Session, postID int) int {
	var result LockDBTemplate
	session.DB("luogulo").C("discussLock").Find(bson.M{"postid": postID}).One(&result)
	return result.LockStatus
}
