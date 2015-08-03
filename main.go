// main.go
package main

import (
	"log"
	"time"
)

// chan
var c chan int

/**
  start thread
  **/
func startKeepAlive(url UrlObj) {
	interval := time.Duration(url.Interval)
	timer := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-timer.C:
			isDie := KeepAlive(url.Url, url.Timeout, 2)
			if isDie {
				var to string
				for index, receiver := range ConfigObj.Receivers {
					to += receiver
					if index != len(ConfigObj.Receivers)-1 {
						to += ";"
					}
				}
				var subject string
				if len(url.Titile) <= 0 {
					subject = "Http keep alive"
				} else {
					subject = url.Titile
				}
				body := "<html><body><h3> " + url.Url + " is die </h3></body></html>"
				sender := ConfigObj.EmailSender
				err := SendMail(sender.Username, sender.Password, sender.Server, to, subject, body, "html")
				if err != nil {
					log.Println("[error] send mail to", to, err.Error())
				} else {
					log.Println("[success] send mail to", to)
				}
			}
		}
	}
	c <- 1
}

/**
  main function
  **/
func main() {
	c = make(chan int)
	for _, urlObj := range ConfigObj.Urls {
		// goroutine
		go startKeepAlive(urlObj)
	}
	<-c
}
