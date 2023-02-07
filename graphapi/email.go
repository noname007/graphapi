package graphapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Conf struct {
	TenantId     string
	ClientId     string
	ClientSecret string
}

type Client struct {
	c     *resty.Client
	Conf  Conf
	token string
}

func (c *Client) fetchToken() {
	body := map[string]string{
		"grant_type":    "client_credentials",
		"scope":         "https://microsoftgraph.chinacloudapi.cn/.default",
		"client_id":     c.Conf.ClientId,
		"client_secret": c.Conf.ClientSecret,
	}
	url := fmt.Sprintf("https://login.chinacloudapi.cn/%s/oauth2/v2.0/token", c.Conf.TenantId)
	post, err := c.c.R().SetFormData(body).Post(url)

	fmt.Println(post, err)
	var resp map[string]any
	err = json.Unmarshal(post.Body(), &resp)
	token := resp["access_token"].(string)
	fmt.Println(err, resp["access_token"])
	c.token = token
}

func (c *Client) FetchEmail(ctx context.Context, upn, subject string) {
	url := fmt.Sprintf("https://microsoftgraph.chinacloudapi.cn/v1.0/users/%s/messages?$filter=subject eq '%s'", upn, subject)
	response, err := c.c.R().SetAuthToken(c.token).Delete(url)
	fmt.Println(response, err)
}

func (c *Client) DeleteEmail(ctx context.Context, upn, msgId string) {
	//url = fmt.Sprintf("https://microsoftgraph.chinacloudapi.cn/v1.0/users/%s/messages?$select=sender,subject", upn)
	//get, err := c.c.R().SetAuthToken(c.token).Get(url)
	//fmt.Println(get, err)
	url := fmt.Sprintf("https://microsoftgraph.chinacloudapi.cn/v1.0/users/%s/messages/%s", upn, msgId)
	response, err := c.c.R().SetAuthToken(c.token).Delete(url)
	fmt.Println(response, err)
}

func NewClient(conf Conf) *Client {
	r := resty.New()

	r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	r.EnableTrace().SetDebug(true)

	c := &Client{
		c:     r,
		Conf:  conf,
		token: "",
	}

	c.fetchToken()
	go func() {
		//TODO update token period
		//ref1 man
		//ref2 https://learnku.com/articles/23578/the-difference-between-go-timer-and-ticker
		//t := time.NewTimer(10 * time.Second)
		t := time.NewTicker(120 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				c.fetchToken()
			}
		}
	}()

	return c
}
