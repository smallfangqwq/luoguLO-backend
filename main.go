package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

var Config Configurations
var Session *mgo.Session
var err error

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("[ERROR] ", err, "\n")
			time.Sleep(1 * time.Second)
			fmt.Print("[INFO] Restart. \n")
			main()
		}
	}()
	if _, err := toml.DecodeFile(os.Args[1], &Config); err != nil {
		panic(err)
	}
	if Session, err = mgo.Dial(Config.Database.URL); err != nil {
		panic(err)
	} else {
		defer Session.Close()
	}
	go autoSaveMain()
	go webMain()
	for {
	}
}
