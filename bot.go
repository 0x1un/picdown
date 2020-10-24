package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	BaseUrl  = `https://oapi.dingtalk.com/robot/send?access_token=`
	FileName = "log/chatbot.log"
)

type Message struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func Send(tokens, atUsers []string, notifyAll bool, text, title string) {
	logFile, _ := os.OpenFile(FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	Log := log.New(logFile, "[Info]", log.Ldate|log.Ltime) // log.Ldate|log.Ltime|log.Lshortfile
	Log.Println("开始发送消息!")
	msg := &Message{}
	msg.MsgType = "markdown"
	msg.Markdown.Title = title
	msg.Markdown.Text = text
	msg.At.AtMobiles = atUsers
	msg.At.IsAtAll = notifyAll
	msgs, err := json.Marshal(msg)
	if err != nil {
		Log.Fatal(err)
	}
	fmt.Println(string(msgs))
	for _, tk := range tokens {
		fillMsgAndSent(tk, msgs, Log)
	}
}

//发送消息到钉钉
func fillMsgAndSent(token string, msg []byte, Log *log.Logger) {
	reader := bytes.NewReader(msg)
	resp := Post(BaseUrl+token, reader)
	Log.SetPrefix("[Info]")
	Log.Printf("消息发送完成,服务器返回内容：%s", string(resp))
}

func Post(url string, reader *bytes.Reader) []byte {
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	return respBytes
}
