package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("[ERROR] ", err, "\n")
			time.Sleep(1 * time.Second)
			fmt.Print("[INFO] Restart. \n")
			main()
		}
	}()
	var config Configurations
	if _, err := toml.DecodeFile(os.Args[1], &config); err != nil {
		panic(err)
	}
	if session, err := mgo.Dial(config.Database.URL); err != nil {
		panic(err)
	} else {
		defer session.Close()
		go AutoSave(session, config)
	}
	for {
	}
}
