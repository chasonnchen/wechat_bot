package xiaojphone

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/chasonnchen/wechat_bot/lib/http"
)

// 查询手机靓号价格等信息，数据来源 http://xx.xiaoj.com/
// api中sign 加密方式，  sign=md5(params + data + $2019%^zjlh)

var (
	sinaClient = &Client{}
)

type Client struct {
}

func NewClient() *Client {
	return sinaClient
}

type Info struct {
	Name    string `json:"name"`
	Content int64  `json:"content"`
}
type Data struct {
	// 归属地
	Regional string `json:"regional"`
	// 运营商
	Operator string `json:"operator"`
	// 详细介绍
	Detail string  `json:"detail"`
	Info   []*Info `json:"info"`
}

type ResponseBody struct {
	Flag string `json:"flag"`
	Data *Data  `json:"data"`
}

func (c *Client) getSign(strParams string, strData string) string {
	data := []byte(strParams + strData + "$2019%^zjlh")
	h := md5.New()
	h.Write(data)
	s := hex.EncodeToString(h.Sum(nil))
	return s
}

func (c *Client) GetInfo(phone string) string {
	strUri := "http://api.xiaoj.com/api/mobile/detail"
	strDate := time.Now().Format("20060102")
	strParams := `{"number":` + phone + `,"domain":"xx.xiaoj.com"}`
	sign := c.getSign(strParams, strDate)

	postMap := make(map[string]string, 3)
	postMap["date"] = strDate
	postMap["param"] = strParams
	postMap["sign"] = sign

	log.Printf("xiaoj req map is %#v", postMap)
	httpClient := http.NewHttpClient("", 2*time.Second)
	resData, err := httpClient.PostForm(strUri, postMap, http.PostOptions{})
	if err != nil {
		log.Printf("xiaoj request err. err = %v", err)
		return "亲，靓号好像不对哦，没有查到呢~\n 您检查下重新发我吧~"
	}

	resObj := new(ResponseBody)
	err = json.Unmarshal([]byte(resData), &resObj)
	if err != nil {
		log.Printf("xiaoj json decode err  is %#v", err)
		return ""
	}

	msg := "手机号: " + phone + "\n"
	msg = msg + "归属地: " + resObj.Data.Regional + "\n"
	msg = msg + "运营商: " + resObj.Data.Operator + "\n"
	msg = msg + "号详情: " + resObj.Data.Detail + "\n"
	msg = msg + "--------------------\n"
	for _, item := range resObj.Data.Info {
		msg = msg + item.Name + ": " + strconv.FormatInt(item.Content, 10) + "\n"
	}

	return msg

	return string(resData)
}
