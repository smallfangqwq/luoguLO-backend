package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SaveAlone SaveAlone按顺序存！！！111111
func SaveAlone(session *mgo.Session, PostID int, FromPage int, EndPage int) { // [From Page, End Page]
	// 检查是否已经存在
	var discuss DBDiscussTemplate
	discussCount, err := session.DB("luogulo").C("discuss").Find(bson.M{"postid": PostID}).Count()
	if err != nil {
		panic(err)
	}
	if discussCount == 0 {
		nowThings := ChangeDiscussSomePageToDBDiscussTemlate(PostID, 1, EndPage)
		err = session.DB("luogulo").C("discuss").Insert(&nowThings)
		if err != nil {
			fmt.Print("[Save ERROR] ERROR1. LOG:", err, "\n")
		}
	} else {
		err = session.DB("luogulo").C("discuss").Find(bson.M{"postid": PostID}).One(&discuss)
		if err != nil {
			panic(err)
		}
		// 先看看爬到的最后一条的发布时间
		var lastTime int64
		if discuss.Count > 0 {
			lastTime = discuss.Comment[discuss.Count-1].SendTime
		} else {
			lastTime = 0
		}
		// 分析现在的帖子
		nowThings := ChangeDiscussSomePageToDBDiscussTemlate(PostID, FromPage, EndPage)
		// 将新评论整理
		lens := nowThings.Count
		NewDiscuss := discuss
		for i := 0; i < lens; i++ {
			if nowThings.Comment[i].SendTime > lastTime {
				NewDiscuss.Comment = append(NewDiscuss.Comment, nowThings.Comment[i])
				NewDiscuss.Count++
			} else if nowThings.Comment[i].SendTime == lastTime {
				flag := false
				for j := 0; j < discuss.Count; j++ {
					if discuss.Comment[j].Content == nowThings.Comment[j].Content {
						flag = true
						break
					}
				}
				if !flag {
					NewDiscuss.Comment = append(NewDiscuss.Comment, nowThings.Comment[i])
					NewDiscuss.Count++
				}
			}
		}
		err = session.DB("luogulo").C("discuss").Update(&discuss, &NewDiscuss)
		if err != nil {
			panic(err)
		}
	}
}

func ChangeDiscussSomePageToDBDiscussTemlate(PostID int, FromPage int, EndPage int) (result DBDiscussTemplate) {
	url := ""
	for i := FromPage; i <= EndPage; i++ {
		url = "https://www.luogu.com.cn/discuss/" + strconv.Itoa(PostID) + "?page=" + strconv.Itoa(i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Set("Cookie", Config.Request.Cookie)
		req.Header.Set("Host", "www.luogu.com.cn")
		req.Header.Set("Referer", "https://www.luogu.com.cn")
		req.Header.Set("User-Agent", Config.Request.UserAgent)
		client := &http.Client{Timeout: time.Second * 15}
		resp, err := client.Do(req)
		if err != nil {
			break
		}
		newDoc, _ := goquery.NewDocumentFromReader(resp.Body)
		nowCounts := result.Count
		newDoc.Find(".am-comment-meta").Each(func(i int, selection *goquery.Selection) {
			texts := selection.Find("a").First().Text()
			if i == 0 {
				if i == FromPage {
					oldT := selection.Text()
					regR, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
					sendTimes, _ := time.Parse("2006-01-02 15:04", regR.FindString(oldT))
					result.SendTime = sendTimes.Unix()
					result.Author = texts
					AuthorId, _ := selection.Find("a").First().Attr("href")
					AuthorId = strings.Trim(AuthorId, "/user")
					result.AuthorID = AuthorId
				}
				return
			}
			result.Count++
			oldT := selection.Text()
			regR, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
			var newComment DBComment
			newComment.Author = texts
			sendTimes, _ := time.Parse("2006-01-02 15:04", regR.FindString(oldT))
			newComment.SendTime = sendTimes.Unix()
			AuthorId, _ := selection.Find("a").First().Attr("href")
			AuthorId = strings.Trim(AuthorId, "/user")
			newComment.AuthorId = AuthorId
			result.Comment = append(result.Comment, newComment)
		})
		newDoc.Find(".am-comment-bd").Each(func(i int, selection *goquery.Selection) {
			htmls, _ := selection.Html()
			if i == 0 {
				return
			}
			result.Comment[i-1+nowCounts].Content = htmls
		})
		if nowCounts == result.Count {
			resp.Body.Close()
			break
		}
		resp.Body.Close()
	}
	return
}
