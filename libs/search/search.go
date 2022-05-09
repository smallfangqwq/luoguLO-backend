package search

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBComment struct {
	SendTime int64
	Author   string
	AuthorId string
	Content  string
}

type DBDiscussTemplate struct {
	PostID   int
	Author   string
	AuthorID string
	SendTime int64
	Title    string
	Describe string
	Count    int
	Comment  []DBComment
}

func FindPostByTime (session *mgo.Session, time int64) []DBDiscussTemplate {
	var answer []DBDiscussTemplate
	session.DB("luogulo").C("discuss").Find(bson.M{"sendtime": time}).All(&answer)
	return answer
}