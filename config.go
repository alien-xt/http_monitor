// config
package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	EmailSender EmailSender `json:"emailSender"`
	Urls        []UrlObj    `json:"urls"`
	Receivers   []string    `json:"receivers"`
}

type UrlObj struct {
	Url      string `json:"url"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	Titile   string `json:"title"`
}

type EmailSender struct {
	Server   string `json:"smtpServer"`
	Username string `json:"userName"`
	Password string `json:"passWord"`
}

var ConfigObj Config
var ConfiFilePath string = "config.json"

/**
  init function
  **/
func init() {
	argNum := len(os.Args)
	if argNum > 1 {
		for index, param := range os.Args {
			if param == "--config" {
				ConfiFilePath = os.Args[index+1]
			}
		}
	}
	data := ReadFile(ConfiFilePath)
	if data == nil {
		log.Fatalln("[error] read config file", 2)
		return
	}
	err := json.Unmarshal(data, &ConfigObj)
	if err != nil {
		log.Fatalln("[error] format config file", 3, err.Error())
		return
	}
	log.Println("[success] load config", ConfigObj)
}
