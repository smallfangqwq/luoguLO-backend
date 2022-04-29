package user

import (
	"awesomeProject/src/database"
	"go.mongodb.org/mongo-driver/bson"
)

func checkToken(token string, userID int) (map[string]interface{}, int) {
	findData := bson.M{
		"userID": userID,
		"token":  token,
		"status": "online",
	}
	dataList, cnt := database.SQLFindDataGroupByConditionWithNew("token", findData)
	if cnt < 1 {
		return nil, -1
	}
	if dataList[1]["ver"].(int) < 2 {
		return (dataList[1]["limit"]).(map[string]interface{}), 0
	}
	return nil, -2
}

func genToken(lens int) string {

	return ""
}

func createToken(userID int, limit interface{}) (string, bool) {
	var newData = make(map[string]interface{})
	newData["userID"] = userID
	newData["ver"] = 1
	var newToken = genToken(16)
	newData["token"] = newToken
	newData["limit"] = limit
	newData["status"] = "online"
	status := database.SQLAddDataWithNew("token", newData)
	return newToken, status
}

func deleteToken(token string) bool {
	findData := bson.M{
		"token": token,
	}
	_, cnt := database.SQLFindDataGroupByConditionWithNew("token", findData)
	if cnt < 1 {
		return false
	}
	//TODO: database update data.
	//TODO: database.SQLUpdateDataWithNew("token", newData)
	return true
}
