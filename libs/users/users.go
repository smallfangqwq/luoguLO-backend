package users

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBUserStruct struct {
	userName string
	userID int
	userLUOGUID int
	userAccess int
}

func GetUserDataByUserID(session *mgo.Session, userID int) (UserData DBUserStruct, err error) {
	err = session.DB("luogulo").C("user").Find(bson.M{"userID": userID}).One(&UserData)
	if err != nil {
		return
	}
	return 
}

func CreateANewUser(session *mgo.Session, user DBUserStruct) {
	session.DB("luogulo").C("user").Insert(user)
}

func RegisterUser(session *mgo.Session, userName string, userLUOGUID int, userAccess int) (error) {
	// get new user ID
	var UserInfo DBUserStruct
	var err error
	UserInfo.userID, err = session.DB("luogulo").C("user").Count()
	if err != nil {
		return err
	}
	UserInfo.userID ++
	UserInfo.userName = userName
	UserInfo.userLUOGUID = userLUOGUID
	UserInfo.userAccess = userAccess
	CreateANewUser(session, UserInfo)
	return nil
}