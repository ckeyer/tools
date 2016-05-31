package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

const (
	user_email    = "me@luv.gift"
	user_password = "zhihuData"
)

var (
	zhihuIndexURL = "http://www.zhihu.com/#signin"
	zhihuLoginURL = "http://www.zhihu.com/login/email"
	zhihuHeader   = map[string]string{
		"Host":            "www.zhihu.com",
		"User-Agent":      "Ckeyer/1.0 (Macintosh; Intel Linux 4.10; rv:16.3) Firefox/41.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "zh-CN,en-US;q=0.7,en;q=0.3' --compressed",
		"Connection":      "keep-alive",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
	}
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type Msg struct {
	R       int   `json:"r"`
	ErrCode int64 `json:"errcode"`
	Data    struct {
		Account string `json:"account"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type CookieMap map[string]*http.Cookie

func (c CookieMap) GetValue(key string) string {
	if v, ok := c[key]; ok {
		return v.Value
	}
	return ""
}

type Client struct {
	req     *http.Request
	cookies CookieMap
}

func NewClient(email, password string) *Client {
	req, err := http.NewRequest("GET", zhihuIndexURL, nil)
	if err != nil {
		log.Fatalln(err)
		// return
	}
	for k, v := range zhihuHeader {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		// return
	}

	cookies := make(CookieMap)
	for k, v := range resp.Cookies() {
		req.AddCookie(v)
		cookies[v.Name] = v
		log.Debugf("%v: %+v", k, v)
	}

	params := make(url.Values)
	params.Add("_xsrf", cookies.GetValue("_xsrf"))
	params.Add("password", password)
	params.Add("email", email)
	params.Add("remember_me", "true")
	params.Add("captcha_type", "cn")

	req.Method = "POST"
	req.URL, _ = url.Parse(zhihuLoginURL)
	req.PostForm = params

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		// return
	}

	var msg Msg
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("return msg: %+v", msg)

	return &Client{req, cookies}
}

func main() {

	c := NewClient(user_email, user_password)
	_ = c
}
