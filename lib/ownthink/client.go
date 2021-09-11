package ownthink

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chasonnchen/wechat_bot/lib/http"
)

// 智能聊天机器人，平台地址https://www.ownthink.com/docs/

var (
	client = &Client{}
)

type Client struct {
}

type Info struct {
	Text      string    `json:"text"`
	Heuristic []*string `json:"heuristic"`
}

type Data struct {
	Type int32 `json:"type"`
	Info *Info `json:"info"`
}

type ResponseBody struct {
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

func NewClient() *Client {
	return client
}

func (c *Client) Ask(id string, text string) string {
	strUri := fmt.Sprintf("https://api.ownthink.com/bot?appid=e97d564f187cba33d3e91b15ee91d285&userid=%s&spoken=%s", id, text)

	log.Printf("ownthink req uri is %s", strUri)
	httpClient := http.NewHttpClient("", 10*time.Second)
	resData, err := httpClient.Get(strUri, nil, http.GetOptions{})
	if err != nil {
		log.Printf("ownthink request err. err = %v", err)
		return ""
	}

	resObj := new(ResponseBody)
	err = json.Unmarshal(resData, &resObj)
	if err != nil {
		log.Printf("json decode err  is %#v", err)
	}

	log.Printf("strUri response data is %#v", resObj)

	return resObj.Data.Info.Text
}
