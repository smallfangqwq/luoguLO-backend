package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var timeInterval interface{}
var timeOlder interface{}

var fetchCount int64 = 0

type AutoSaveFetchStruct struct {
	Status int `json:"status"`
	Data   struct {
		Count  int `json:"count"`
		Result []struct {
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
		} `json:"result"`
	} `json:"data"`
}

func JSONToMap(str []byte) AutoSaveFetchStruct {
	var tempMap AutoSaveFetchStruct
	err := json.Unmarshal(str, &tempMap)

	if err != nil {
		fmt.Printf("[ERROR] JSON DEAL ERROR WITH AUTOSAVE, LOG:\n", err)
		return tempMap
	}

	return tempMap
}

type DBDiscussTemplate struct {
	PostID   int
	Author   string
	title    string
	describe string
	count    int
	Comment  []struct {
		SendTime int64
		Author   string
		Content  string
	}
}

func ChangeDiscussToDBDiscussTemlate(PostID int) (result DBDiscussTemplate) {
	//这个要爬多页面的...离谱
	//想想都头疼...
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
	doc.Find("article").Each(func(i int, selection *goquery.Selection) {
		fmt.Println("i", i, "select text", selection)
	})
	return
}

func SaveNewDiscuss(PostID int) {
	// 检查是否已经存在
	session, err := mgo.Dial("mongodb://root:rtpwd@localhost:62232")
	var discuss DBDiscussTemplate
	var discussCount int
	discussCount, err = session.DB("luogulo").C("discuss").Find(bson.M{"id": PostID}).Count()
	if err != nil {
		fmt.Print("[Save ERROR] Can`t check exist. LOG:", err)
		return
	}
	if discussCount == 0 {
		// 爬全部
		// 分析帖子
	} else {
		err = session.DB("luogulo").C("discuss").Find(bson.M{"id": PostID}).One(&discuss)
		if err != nil {
			fmt.Print("[Save ERROR] Can`t read discuss information before. LOG:", err)
			return
		}
		// 先看看爬到的最后一条的发布时间
		lastTime := discuss.Comment[discuss.count-1].SendTime
		// 分析现在的帖子
		nowThings := ChangeDiscussToDBDiscussTemlate(PostID)
		// 将新评论整理
		lens := nowThings.count
		NewDiscuss := discuss
		for i := 0; i < lens; i++ {
			if nowThings.Comment[i].SendTime > lastTime {
				NewDiscuss.Comment[NewDiscuss.count] = nowThings.Comment[i]
				NewDiscuss.count++
			}
		}
		// 将内容update.
		newData, err := bson.Marshal(&NewDiscuss)
		if err != nil {
			fmt.Print("[Save ERROR] ERROR. LOG:", err)
		}
		oldData, err := bson.Marshal(&discuss)
		session.DB("luogulo").C("discuss").Update(oldData, newData)
	}
}

func AutoSave() {
	runtime.Gosched()
	fmt.Printf("[Info] AutoSave Tool has been started.\n")
	//time.Sleep(2 * time.Second)
	for true {
		// runtime.Gosched()
		//fmt.Printf("[Info] AutoSave Tool are fetching from Luogu now. \n") 隐藏本条并修改语法。
		url := "https://www.luogu.com.cn/api/discuss?forum=relevantaffairs&page=1"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("[ERROR] AutoSave Tool can`t get Luogu discuss now. We change the interval to 120 seconds.\n")
			time.Sleep(120 * time.Second)
			continue
		}
		req.Header.Set("Cookie", "")
		req.Header.Set("Host", "www.luogu.com.cn")
		req.Header.Set("Referer", "https://www.luogu.com.cn")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
		client := &http.Client{Timeout: time.Second * 15}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[ERROR] AutoSave Tool can`t get Luogu discuss now. We change the interval to 120 seconds.\n")
			fmt.Print("[ERROR]Error reading response. ", err)
			time.Sleep(120 * time.Second)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		dataThings := JSONToMap(body)
		lenOfDataResult := len(dataThings.Data.Result)
		fetchCount++
		for i := 0; i < lenOfDataResult; i++ {
			var discuss = dataThings.Data.Result[i]
			if discuss.Top > 2 {
				// 忽略置顶
				continue
			}
			SaveNewDiscuss(dataThings.Data.Result[i].PostID)
		}
		time.Sleep(timeInterval.(time.Duration))
	}
}

func main() {
	timeInterval = 5 * 1000 * time.Millisecond
	timeOlder = timeOlder
	go AutoSave()
	for true {
		var command string
		fmt.Scanln(&command)
		fmt.Printf("[Info]Deal with " + command + "\n")
		if command == "changeTime" || command == "ct" {
			fmt.Printf("[AutoSave] How long do you want to set(millonsecond)?\n")
			var newTime int64
			fmt.Scanln(&newTime)
			timeInterval = time.Duration(newTime) * time.Millisecond
			timeOlder = timeInterval
			fmt.Printf("[AutoSave] Settings done!\n")
		}
		if command == "countf" || command == "ft" {
			fmt.Printf("[AutoSave] We fetch %d time(s)\n", fetchCount)
		}
		if command == "debugD" || command == "dd" {
			fmt.Printf("[Discuss] ID?\n")
			var discussID int
			fmt.Scanln(&discussID)
			ChangeDiscussToDBDiscussTemlate(discussID)
		}
	}
}
