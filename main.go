package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

type LegacyPost struct {
	PostID int    `json:"PostID"`
	Title  string `json:"Title"`
	Author struct {
		Instance string `json:"_instance"`
	} `json:"Author"`
	Forum struct {
		ForumID      int    `json:"ForumID"`
		Name         string `json:"Name"`
		InternalName string `json:"InternalName"`
		Instance     string `json:"_instance"`
	} `json:"Forum"`
	Top         int  `json:"Top"`
	SubmitTime  int  `json:"SubmitTime"`
	IsValid     bool `json:"isValid"`
	LatestReply struct {
		Author struct {
			Instance string `json:"_instance"`
		} `json:"Author"`
		ReplyTime int    `json:"ReplyTime"`
		Content   string `json:"Content"`
		Instance  string `json:"_instance"`
	} `json:"LatestReply"`
	RepliesCount int    `json:"RepliesCount"`
	Instance     string `json:"_instance"`
}

type LegacyPostList struct {
	Status int `json:"status"`
	Data   struct {
		Count  int          `json:"count"`
		Result []LegacyPost `json:"result"`
	} `json:"data"`
}

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

func ChangeDiscussToDBDiscussTemlate(PostID int) (result DBDiscussTemplate) {
	//Complete !
	url := "https://www.luogu.com.cn/discuss/" + strconv.Itoa(PostID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR] Tool can`t get Luogu discuss now.\n")
		time.Sleep(120 * time.Second)
		return
	}
	req.Header.Set("Cookie", "UM_distinctid=17d89339530bd4-0df461f9ef6091-1f396452-13c680-17d89339531ceb; login_referer=https%3A%2F%2Fwww.luogu.com.cn%2F; __client_id=ae4f59efbd21087f9cb79c186e8d4d91044e0db9; _uid=99640; CNZZDATA5476811=cnzz_eid%3D613104886-1624186548-%26ntime%3D1651295299")
	req.Header.Set("Host", "www.luogu.com.cn")
	req.Header.Set("Referer", "https://www.luogu.com.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 15}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Tool can`t get Luogu discuss now.\n")
		fmt.Print("[ERROR]Error reading response. ", err)
		time.Sleep(120 * time.Second)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Tool can`t get Luogu discuss now.\n")
		fmt.Print("[ERROR]Error reading response. ", err)
		time.Sleep(120 * time.Second)
		return
	}
	result.Count = 0
	titles := doc.Find("h1").First().Text()
	result.Title = titles
	result.PostID = PostID

	result.Count = 0
	//fmt.Printf(titles)
	// 获取每条评论的发布时间和内容以及整个帖子内容
	doc.Find(".am-comment-meta").Each(func(i int, selection *goquery.Selection) {
		texts := selection.Find("a").First().Text()
		if i == 0 {
			oldT := selection.Text()
			regR, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
			sendTimes, _ := time.Parse("2006-01-02 15:04", regR.FindString(oldT))
			result.SendTime = sendTimes.Unix()
			result.Author = texts
			AuthorId, _ := selection.Find("a").First().Attr("href")
			AuthorId = strings.Trim(AuthorId, "/user/")
			result.AuthorID = AuthorId
			return
		}
		result.Count++
		oldT := selection.Text()
		regR, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
		//result.Comment[i].Author = texts
		var newComment DBComment
		newComment.Author = texts
		AuthorId, _ := selection.Find("a").First().Attr("href")
		AuthorId = strings.Trim(AuthorId, "/user/")
		result.AuthorID = AuthorId
		sendTimes, _ := time.Parse("2006-01-02 15:04", regR.FindString(oldT))
		newComment.SendTime = sendTimes.Unix()
		result.Comment = append(result.Comment, newComment)
		//fmt.Println("i", i, "select text", texts)
	})
	//每条内容内容获取和标题内容
	doc.Find(".am-comment-bd").Each(func(i int, selection *goquery.Selection) {
		htmls, _ := selection.Html()
		if i == 0 {
			result.Describe = htmls
			return
		}
		result.Comment[i-1].Content = htmls
		//	fmt.Println("i", i, "select text", htmls)
	})
	// 多页面
	// TODO: 页面过多分段处理问题
	for i := 2; true; i++ {
		url = "https://www.luogu.com.cn/discuss/" + strconv.Itoa(PostID) + "?page=" + strconv.Itoa(i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			//fmt.Printf("[ERROR] Tool can`t get Luogu discuss now.\n")
			break
		}
		req.Header.Set("Cookie", "UM_distinctid=17d89339530bd4-0df461f9ef6091-1f396452-13c680-17d89339531ceb; login_referer=https%3A%2F%2Fwww.luogu.com.cn%2F; __client_id=ae4f59efbd21087f9cb79c186e8d4d91044e0db9; _uid=99640; CNZZDATA5476811=cnzz_eid%3D613104886-1624186548-%26ntime%3D1651295299")
		req.Header.Set("Host", "www.luogu.com.cn")
		req.Header.Set("Referer", "https://www.luogu.com.cn")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
		client := &http.Client{Timeout: time.Second * 15}
		resp, err := client.Do(req)
		if err != nil {
			//	fmt.Printf("[ERROR] Tool can`t get Luogu discuss now.\n")
			//	fmt.Print("[ERROR]Error reading response. ", err)
			break
		}
		defer resp.Body.Close()
		newDoc, _ := goquery.NewDocumentFromReader(resp.Body)
		nowCounts := result.Count
		newDoc.Find(".am-comment-meta").Each(func(i int, selection *goquery.Selection) {
			texts := selection.Find("a").First().Text()
			if i == 0 {
				return
			}
			result.Count++
			oldT := selection.Text()
			regR, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
			//result.Comment[i].Author = texts
			var newComment DBComment
			newComment.Author = texts
			sendTimes, _ := time.Parse("2006-01-02 15:04", regR.FindString(oldT))
			newComment.SendTime = sendTimes.Unix()
			result.Comment = append(result.Comment, newComment)
			//fmt.Println("i", i, "select text", texts)
		})
		//每条内容内容获取和标题内容
		newDoc.Find(".am-comment-bd").Each(func(i int, selection *goquery.Selection) {
			htmls, _ := selection.Html()
			if i == 0 {
				return
			}
			result.Comment[i-1+nowCounts].Content = htmls
			//	fmt.Println("i", i, "select text", htmls)
		})
		if nowCounts == result.Count { // 到头啦
			break
		}
	}
	//	fmt.Print(result)
	return
}

func SaveNewDiscuss(session *mgo.Session, PostID int) {
	// 检查是否已经存在
	var discuss DBDiscussTemplate
	discussCount, err := session.DB("luogulo").C("discuss").Find(bson.M{"postid": PostID}).Count()
	if err != nil {
		fmt.Print("[Save ERROR] Can`t check exist. LOG:", err)
		return
	}
	if discussCount == 0 {
		// 爬全部
		//	fmt.Printf("HERE")
		nowThings := ChangeDiscussToDBDiscussTemlate(PostID)
		// 分析帖子
		//nowThingsDB, err := bson.Marshal(&nowThings)
		//if err != nil {
		//	fmt.Print("[Save ERROR] ERROR. LOG:", err)
		//}
		err = session.DB("luogulo").C("discuss").Insert(&nowThings)
		if err != nil {
			fmt.Print("[Save ERROR] ERROR1. LOG:", err, "\n")
		}
	} else {
		err = session.DB("luogulo").C("discuss").Find(bson.M{"postid": PostID}).One(&discuss)
		if err != nil {
			fmt.Print("[Save ERROR] Can`t read discuss information before. LOG:", err, "\n")
			return
		}
		// 先看看爬到的最后一条的发布时间
		var lastTime int64
		if discuss.Count > 0 {
			lastTime = discuss.Comment[discuss.Count-1].SendTime
		} else {
			lastTime = 0
		}

		// 分析现在的帖子
		nowThings := ChangeDiscussToDBDiscussTemlate(PostID)
		// 将新评论整理
		lens := nowThings.Count
		NewDiscuss := discuss
		for i := 0; i < lens; i++ {
			if nowThings.Comment[i].SendTime > lastTime {
				NewDiscuss.Comment = append(NewDiscuss.Comment, nowThings.Comment[i])
				NewDiscuss.Count++
			} else if nowThings.Comment[i].SendTime == lastTime { // 可爱的洛谷竟然只到分钟，显然有可能遇到时间问题
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
		// 将内容update.
		err = session.DB("luogulo").C("discuss").Update(&discuss, &NewDiscuss)
		if err != nil {
			fmt.Print("[Save ERROR] ERROR2. LOG:", err, "\n")
		}
	}
}

func AutoSave(session *mgo.Session) {
	url := "https://www.luogu.com.cn/api/discuss?forum=relevantaffairs&page=1"
	fmt.Println("Listing", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}

	var result LegacyPostList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	} else if result.Status != http.StatusOK {
		panic(result.Status)
	}

	lenOfDataResult := len(result.Data.Result)
	for i := 0; i < lenOfDataResult; i++ {
		fmt.Printf("\tSaving %d...\n", result.Data.Result[i].PostID)
		SaveNewDiscuss(session, result.Data.Result[i].PostID)
	}

	fmt.Println("Fetched", url)
}

func main() {
	if session, err := mgo.Dial("mongodb://localhost:27017"); err != nil {
		panic(err)
	} else {
		defer session.Close()
		AutoSave(session)
	}
}
