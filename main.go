package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2"
)

type DatabaseConfigurations struct {
	URL string
}

type RequestConfigurations struct {
	UserAgent string `toml:"user_agent"`
	Cookie    string
}
type HttpConfigurations struct {
	Port         int
	AccessOrigin string
}

type FunctionConfigurations struct {
	Autosave bool
	Web      bool
}

type Configurations struct {
	Database     DatabaseConfigurations
	Request      RequestConfigurations
	Http         HttpConfigurations
	Function     FunctionConfigurations `toml:"functions"`
	Target       string
	TimeInterval int `toml:"time_interval"`
}

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
	var wg sync.WaitGroup
	if Config.Function.Autosave {
		wg.Add(1)
		go autoSaveMain()
	}
	if Config.Function.Web {
		wg.Add(1)
		go webMain()
	}
	wg.Wait()
}
