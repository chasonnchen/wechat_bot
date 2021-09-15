package baidu

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chasonnchen/wechat_bot/lib/http"
)

type AccessToken struct {
	Token       string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ExpiresTime int64
}

func GenAccessToken(ak string, sk string) (accessToken AccessToken, err error) {
	strUri := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", ak, sk)
	log.Printf("access token get uri is %s", strUri)

	httpClient := http.NewHttpClient("", 2*time.Second)
	resData, err := httpClient.Get(strUri, nil, http.GetOptions{})
	if err != nil {
		log.Printf("access token request err. err = %v", err)
		return AccessToken{}, fmt.Errorf("access token request err. err = %v", err)
	}

	resObj := AccessToken{}
	err = json.Unmarshal(resData, &resObj)
	if err != nil {
		log.Printf("json decode err is %#v", err)
	}
	resObj.ExpiresTime = time.Now().Unix() + resObj.ExpiresIn
	log.Printf("access token res is %#v", resObj)

	return resObj, nil
}
