package sina7x24

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chasonnchen/wechat_bot/lib/http"
)

var (
	sinaClient = &Client{}
)

type Client struct {
}

type NewsTag struct {
	Id string
}
type News struct {
	Id         int32      `json:"id"`
	RichText   string     `json:"rich_text"`
	CreateTime string     `json:"create_time"`
	TagList    []*NewsTag `json:"tag"`
}
type Feed struct {
	NewsList []*News `json:"list"`
}
type Data struct {
	Feed *Feed `json:"feed"`
}

type Result struct {
	Data *Data `json:"data"`
}

type ResponseBody struct {
	Result *Result `json:"result"`
}

func NewClient() *Client {
	return sinaClient
}

func (c *Client) GetMsgs(tagId int32, lastMsgId int32) (msgContent string, lastMessageId int32) {
	strUri := fmt.Sprintf("http://zhibo.sina.com.cn/api/zhibo/feed?&page=1&page_size=20&zhibo_id=152&tag_id=%d&dire=f&dpc=1&pagesize=20&id=%d&type=0&_=%d", tagId, lastMsgId, time.Now().UnixNano()/1e6)
	log.Printf("sina get uri is %s", strUri)

	httpClient := http.NewHttpClient("", 2*time.Second)
	resData, err := httpClient.Get(strUri, nil, http.GetOptions{})
	if err != nil {
		log.Printf("sina 7*24 request err. err = %v", err)
		return "", 0
	}

	resObj := new(ResponseBody)
	err = json.Unmarshal(resData, &resObj)
	if err != nil {
		log.Printf("json decode err  is %#v", err)
		return "", 0
	}

	if len(resObj.Result.Data.Feed.NewsList) > 0 {
		for _, news := range resObj.Result.Data.Feed.NewsList {
			if news.Id > lastMessageId {
				lastMessageId = news.Id
			}

			for _, tag := range news.TagList {
				if tag.Id == "9" {
					msgContent = msgContent + news.CreateTime + "\n" + news.RichText + "\n\n"
				}
			}
		}
	}

	msgContent = strings.Trim(msgContent, "\n")

	return msgContent, lastMessageId
}
