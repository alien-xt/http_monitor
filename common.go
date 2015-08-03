// common
package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	RespTime int
	RespCode int
}

/**
  read file
  **/
func ReadFile(fileName string) []byte {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0660)
	if err != nil {
		log.Fatalln("[error] read file", 1, err.Error())
		return nil
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return data
}

/**
  https get
  **/
func HttpsGet(url string) *Response {
	respObj := new(Response)
	start := time.Now()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("[error] send http get", err.Error())
		return nil
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error] read from resp body")
		return nil
	}
	end := time.Now()
	respTimeSecond := int(end.Sub(start).Seconds())
	respObj.RespCode = resp.StatusCode
	respObj.RespTime = respTimeSecond
	return respObj
}

/**
  http get
  **/
func HttpGet(url string) *Response {
	respObj := new(Response)
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[error] send http get", err.Error())
		return nil
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error] read from response body error", err.Error())
		return nil
	}
	end := time.Now()
	respTimeSecond := int(end.Sub(start).Seconds())
	respObj.RespCode = resp.StatusCode
	respObj.RespTime = respTimeSecond
	return respObj
}

/**
  http get
  **/
func SendGet(url string) *Response {
	if strings.HasPrefix(url, "https") {
		return HttpsGet(url)
	} else {
		return HttpGet(url)
	}
}

/**
  http keep alive
  return true:die false:alive
  **/
func KeepAlive(httpUrl string, respTime int, repeat int) bool {
	var die bool
	resp := SendGet(httpUrl)
	if resp == nil {
		return false
	}
	if resp.RespCode != 200 || resp.RespTime > respTime {
		var errorSize int
		for i := 1; i <= repeat; i++ {
			response := SendGet(httpUrl)
			if resp != nil {
				if response.RespCode != 200 || response.RespTime > respTime {
					errorSize++
				}
			}
			log.Println("[info] try again kepp alive", httpUrl)
			time.Sleep(time.Second * 1)
		}
		if errorSize == repeat {
			die = true
			log.Println("[warn]", httpUrl, "respCode is "+
				strconv.Itoa(resp.RespCode)+" , respTime is "+strconv.Itoa(resp.RespTime)+"s")
		}
	} else {
		die = false
		log.Println("[success]", httpUrl, "respCode is "+
			strconv.Itoa(resp.RespCode)+" , respTime is "+strconv.Itoa(resp.RespTime)+"s")
	}
	return die
}

/**
  send email
  **/
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
