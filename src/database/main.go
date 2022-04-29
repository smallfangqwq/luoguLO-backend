package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectSQL() (*mongo.Client, bool) {
	clientOptions := options.Client().ApplyURI("mongodb://root:etqrefLFHERFHELFHekljwfsfeqfe@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return client, true
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return client, true
	}
	return client, false
}

func stopConnect(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func SQLFindDataGroup(client *mongo.Client, collect string) (map[int]map[string]interface{}, int) {
	collections := client.Database("rotriw").Collection(collect)
	var (
		err    error
		cursor *mongo.Cursor
	)
	if cursor, err = collections.Find(context.TODO(), bson.D{}); err != nil {
		return nil, -1
	}
	var results map[int]map[string]interface{}
	var cnt = 0
	for cursor.Next(context.TODO()) {
		var result map[string]interface{}
		var errs error
		errs = cursor.Decode(&result)
		if errs != nil {
			return nil, -1
		}
		cnt++
		results[cnt] = result
	}
	return results, cnt
}

func SQLFindDataGroupByCondition(client *mongo.Client, collect string, confident interface{}) (map[int]map[string]interface{}, int) {
	collections := client.Database("rotriw").Collection(collect)
	var (
		err    error
		cursor *mongo.Cursor
	)
	if cursor, err = collections.Find(context.TODO(), confident); err != nil {
		return nil, -1
	}
	var results map[int]map[string]interface{}
	var cnt = 0
	for cursor.Next(context.TODO()) {
		var result map[string]interface{}
		var errs error
		if errs != nil {
			return nil, -1
		}
		cnt++
		results[cnt] = result
	}
	return results, cnt
}

func SQLAddData(client *mongo.Client, collect string, data interface{}) (*mongo.InsertOneResult, error, bool) {
	collections := client.Database("rotriw").Collection(collect)
	result, err := collections.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err, false
	}
	return result, nil, true
}

func SQLAddManyData(client *mongo.Client, collect string, data []interface{}) (*mongo.InsertManyResult, error, bool) {
	collections := client.Database("rotriw").Collection(collect)
	result, err := collections.InsertMany(context.TODO(), data)
	if err != nil {
		return nil, err, false
	}
	return result, nil, true
}

func SQLUpdateData(client *mongo.Client, collect string, where interface{}, change interface{}) bool {
	collections := client.Database("rotriw").Collection(collect)
	_, err := collections.UpdateOne(context.TODO(), where, change)
	if err != nil {
		return false
	}
	return true
}

func SQLUpdateDataSet(client *mongo.Client, collect string, where interface{}, change interface{}) bool {
	return SQLUpdateData(client, collect, where, bson.M{"$set": change})
}

func SQLAddDataWithNew(collect string, data interface{}) bool {
	client, _ := ConnectSQL()
	_, _, result := SQLAddData(client, collect, data)
	stopConnect(client)
	return result
}

func SQLUpdateDataWithNew(collect string, where interface{}, change interface{}) bool {
	client, err := ConnectSQL()
	if err == true {
		return false
	}
	results := SQLUpdateData(client, collect, where, change)
	stopConnect(client)
	return results
}

func SQLUpdateDataSetWithNew(collect string, where interface{}, change interface{}) bool {
	client, err := ConnectSQL()
	if err == true {
		return false
	}
	results := SQLUpdateDataSet(client, collect, where, change)
	stopConnect(client)
	return results
}

func SQLFindDataGroupWithNew(collect string) (map[int]map[string]interface{}, int) {
	client, err := ConnectSQL()
	if err == true {
		return nil, -1
	}
	results, cnt := SQLFindDataGroup(client, collect)
	stopConnect(client)
	return results, cnt
}

func SQLFindDataGroupByConditionWithNew(collect string, condition interface{}) (map[int]map[string]interface{}, int) {
	client, err := ConnectSQL()
	if err == true {
		return nil, -1
	}
	result, cnt := SQLFindDataGroupByCondition(client, collect, condition)
	stopConnect(client)
	return result, cnt
}

func SQLAddManyDataWithNew(collect string, data []interface{}) bool {
	client, _ := ConnectSQL()
	_, _, result := SQLAddManyData(client, collect, data)
	stopConnect(client)
	return result
}
