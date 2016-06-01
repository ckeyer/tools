package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
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

type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

type Client struct {
	httpClient *http.Client
	cookies    CookieMap
}

func NewClient() *Client {
	httpClient := &http.Client{}
	httpClient.Jar = new(Jar)
	cookies := make(CookieMap)

	return &Client{httpClient, cookies}
}

func (c *Client) Index() (*http.Response, error) {
	resp, err := c.httpClient.Get("https://www.zhihu.com/#signin")
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Cookies() {
		c.cookies[v.Name] = v
	}
	c.httpClient.Jar.SetCookies(resp.Request.URL, resp.Cookies())

	return resp, nil
}

func (c *Client) Login(email, password string) (bool, error) {
	params := make(url.Values)
	params.Add("_xsrf", c.cookies.GetValue("_xsrf"))
	params.Add("password", password)
	params.Add("email", email)
	params.Add("remember_me", "true")
	params.Add("captcha_type", "cn")
	resp, err := c.httpClient.PostForm("https://www.zhihu.com/login/email", params)
	if err != nil {
		return false, err
	}

	var msg Msg
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return false, err
	}
	c.httpClient.Jar.SetCookies(resp.Request.URL, resp.Cookies())

	log.Infof("return msg: %+v", msg)
	if msg.R == 0 {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetPeople(name string) (*goquery.Document, error) {
	peopleURL := "https://www.zhihu.com/people/" + name
	resp, err := c.httpClient.Get(peopleURL)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(resp)
}

func main() {
	c := NewClient()
	if res, err := c.Index(); err != nil {
		log.Errorf("index failed, error: %s", err.Error())
	} else {
		log.Debugf("index status: %S", res.Status)
	}

	if ok, err := c.Login(user_email, user_password); err != nil {
		log.Errorf("login failed", err.Error())
	} else if !ok {
		log.Errorf("login failed, no error message")
	}

	if doc, err := c.GetPeople("ckeyer"); err != nil {
		log.Errorf("get people page failed, error: %s", err.Error())
	}

}
