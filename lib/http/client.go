package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpBase struct {
	client *http.Client
	addr   string
}

func NewHttpClient(addr string, timeout time.Duration) *HttpBase {
	if timeout < 0 {
		timeout = 1 * time.Second
	}

	hb := new(HttpBase)
	hb.addr = addr
	hb.client = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}
	return hb
}

type GetOptions struct {
	Header map[string]string
}

func (this *HttpBase) Get(uri string, data map[string]string, options GetOptions) (body []byte, err error) {
	url := this.addr + uri
	if strings.HasPrefix(uri, "http") {
		url = uri
	}
	req, err := http.NewRequest("GET", url, nil)

	if data != nil {
		q := req.URL.Query()
		for key, value := range data {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	if options.Header != nil {
		for field, value := range options.Header {
			if strings.ToLower(field) == "host" {
				req.Host = value
			}
			req.Header.Set(field, value)
		}
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http code not 200, is %#v", resp.StatusCode)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

type PostOptions struct {
	Header map[string]string
}

func (this *HttpBase) PostJson(uri string, data interface{}, options PostOptions) (body []byte, err error) {
	url := this.addr + uri
	if strings.HasPrefix(uri, "http") {
		url = uri
	}

	jsonStr, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	if options.Header != nil {
		for field, value := range options.Header {
			if strings.ToLower(field) == "host" {
				req.Host = value
			}
			req.Header.Set(field, value)
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       200,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http code not 200")
	}

	log.Printf("post json response string is %s", string(body))

	return body, nil
}

func (this *HttpBase) PostForm(uri string, data map[string]string, options PostOptions) (body []byte, err error) {
	url := this.addr + uri
	if strings.HasPrefix(uri, "http://") {
		url = uri
	}

	// map转成字符串
	dataList := make([]string, 0)
	for key, value := range data {
		dataList = append(dataList, key+"="+value)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(dataList, "&")))
	if err != nil {
		return nil, err
	}

	if options.Header != nil {
		for field, value := range options.Header {
			if strings.ToLower(field) == "host" {
				req.Host = value
			}
			req.Header.Set(field, value)
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http code not 200")
	}
	return body, nil
}
