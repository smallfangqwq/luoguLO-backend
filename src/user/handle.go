package user

import (
	"awesomeProject/src/database"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(userName string, userPwd string, userEmail string, userGroup string) (status int, msg string) {
	status = 0
	msg = "ok."
	findCon := bson.M{
		"username": userName,
	}
	_, cnt := database.SQLFindDataGroupByConditionWithNew("user", findCon)
	if cnt == -1 {
		status = -1
		msg = "SQL error."
		return
	}
	if cnt > 0 {
		status = 1
		msg = "Same user data."
	}
	return
}
