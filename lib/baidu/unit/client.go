package unit

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/chasonnchen/wechat_bot/configs"
	"github.com/chasonnchen/wechat_bot/lib/baidu"
	"github.com/chasonnchen/wechat_bot/lib/http"
)

var (
	client = &Client{}
	once   sync.Once
)

type Client struct {
	Sessions    map[string]Session
	AccessToken baidu.AccessToken
}

type Session struct {
	Settime   int64
	SessionId string
}

func NewClient() *Client {
	client.init()
	return client
}

type Action struct {
	Say string `json:"say"`
}

type Responses struct {
	Status  int32     `json:"status"`
	Msg     string    `json:"msg"`
	Actions []*Action `json:"actions"`
}

type Context struct {
	History []string `json:"SYS_PRESUMED_HIST"`
}

type Result struct {
	Context   *Context     `json:"context"`
	SessionId string       `json:"session_id"`
	Responses []*Responses `json:"responses"`
}

type UnitRespose struct {
	ErrorCode int32   `json:"error_code"`
	ErrorMsg  string  `json:"error_msg"`
	Result    *Result `json:"result"`
}

// 下面开始是request相关结构体
type UnitRequest struct {
	Version   string   `json:"version"`
	ServiceId string   `json:"service_id"`
	LogId     string   `json:"log_id"`
	SessionId string   `json:"session_id"`
	Request   *Request `json:"request"`
}

type Request struct {
	TerminalId string `json:"terminal_id"`
	Query      string `json:"query"`
}

func (c *Client) getSession(contactId string) Session {
	session, ok := c.Sessions[contactId]
	if ok {
		// 15min过期
		if session.Settime-time.Now().Unix() < 900 {
			return session
		}
	}

	// 直接返回一个新的session即可
	return Session{
		SessionId: "",
	}
}

func (c *Client) setSession(contactId string, session Session) {
	newMap := c.Sessions
	newMap[contactId] = session
	c.Sessions = newMap
	log.Printf("set session success")
	return
}

func (c *Client) genLogId(contactId string) string {
	return "webot__" + contactId + "__" + strconv.FormatInt(time.Now().Unix(), 10)
}

func (c *Client) Chat(contactId string, query string) (say string, err error) {
	// 1. 拼接请求体
	// 1.1 拼接URL
	strUri := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/unit/service/v3/chat?access_token=%s", c.AccessToken.Token)
	log.Printf("baidu unit request url is %s", strUri)
	// 1.2 获取session
	session := c.getSession(contactId)
	// 1.3 拼接请求结构体
	unitRequest := UnitRequest{
		Version:   "3.0",
		ServiceId: "S58199",
		LogId:     c.genLogId(contactId),
		SessionId: session.SessionId,
		Request: &Request{
			TerminalId: contactId,
			Query:      query,
		},
	}
	log.Printf("baidu unit request params is %#v", unitRequest)

	httpClient := http.NewHttpClient("https://aip.baidubce.com", 10*time.Second)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	options := http.PostOptions{
		Header: header,
	}
	responseBody, err := httpClient.PostJson(strUri, unitRequest, options)
	unitRes := UnitRespose{}
	err = json.Unmarshal(responseBody, &unitRes)
	if err != nil {
		log.Printf("baidu unit res unmarshal failed, err: %s, value: %s", err.Error(), string(responseBody))
		return "", fmt.Errorf("baidu unit res unmarshal failed, err: %s, value: %s", err.Error(), string(responseBody))
	}

	// set sessionid
	c.setSession(contactId, Session{
		Settime:   time.Now().Unix(),
		SessionId: unitRes.Result.SessionId,
	})

	// 解析结果
	resStr := ""
	for _, responses := range unitRes.Result.Responses {
		for _, actions := range responses.Actions {
			log.Printf("baidu unit response say is %s", actions.Say)
			resStr = actions.Say
			break
		}
		break
	}

	return resStr, nil
}

func (c *Client) init() {
	once.Do(func() {
		c.Sessions = make(map[string]Session)
		c.load()
		go func() {
			for {
				select {
				case <-time.After(time.Second * 3600):
					c.load()
				}
			}
		}()
	})
}

func (c *Client) load() {
	baiduConf := configs.GetConf().Baidu
	c.AccessToken, _ = baidu.GenAccessToken(baiduConf.Ak, baiduConf.Sk)
	// TODO 请空超时session
}
