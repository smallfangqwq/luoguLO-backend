package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var timeInterval interface{}
var timeOlder interface{}

func AutoSave() {
	runtime.Gosched()
	time.Sleep(2 * time.Second)
	for true {
		// runtime.Gosched()
		fmt.Printf("[Info] AutoSave Tool is fetch from Luogu now. \n")
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
			log.Fatal("[ERROR]Error reading response. ", err)
			time.Sleep(120 * time.Second)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(body))
		var dataMap map[string]int
		dataStr := body
		err = json.Unmarshal([]byte(dataStr), &dataMap)
		if err != nil {
			fmt.Printf("[ERROR] AutoSave Tool can`t change luogu discuss to right way.\n")
			log.Fatal("[ERROR] AutoSave ERROR log: ", err)
			time.Sleep(120 * time.Second)
			continue
		}
		time.Sleep(timeInterval.(time.Duration))
	}
}

func main() {
	timeInterval = 5 * 1000 * time.Millisecond
	timeOlder = timeOlder
	go AutoSave()
	fmt.Printf("[Info] AutoSave Tool has been started.\n")
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
	}
}
