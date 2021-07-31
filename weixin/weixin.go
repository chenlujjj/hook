package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

type WechatClient struct {
	url string
	c   *http.Client
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type TextPayload struct {
	MsgType string `json:"msgtype,omitempty"`
	Text    Text   `json:"text,omitempty"`
}
type Markdown struct {
	Content string `json:"content,omitempty"`
}
type MarkdownPayload struct {
	MsgType  string   `json:"msgtype,omitempty"`
	Markdown Markdown `json:"markdown,omitempty"`
}

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewWechatClient(key string) *WechatClient {
	return &WechatClient{
		url: baseURL + "?key=" + key,
		c:   http.DefaultClient, // do not use default client
	}
}

func (wc *WechatClient) SendText(text string) error {
	textPayload := TextPayload{
		MsgType: "text",
		Text: Text{
			Content: text,
		},
	}
	b, _ := json.Marshal(textPayload)
	req, err := http.NewRequest("POST", wc.url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := wc.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var resp Response
	_ = json.Unmarshal(body, &resp)
	if resp.ErrCode != 0 {
		return fmt.Errorf("response errcode: %d, errmsge: %s", resp.ErrCode, resp.ErrMsg)
	}
	return nil
}

func (wc *WechatClient) SendMarkdown(content string) error {
	mdPayload := MarkdownPayload{
		MsgType: "markdown",
		Markdown: Markdown{
			Content: content,
		},
	}
	b, _ := json.Marshal(mdPayload)
	req, err := http.NewRequest("POST", wc.url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := wc.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var resp Response
	_ = json.Unmarshal(body, &resp)
	if resp.ErrCode != 0 {
		return fmt.Errorf("response errcode: %d, errmsge: %s", resp.ErrCode, resp.ErrMsg)
	}
	return nil
}

func (wc *WechatClient) SendImage(text string) error {

	return nil
}

// 发送图文
func (wc *WechatClient) SendNews(text string) error {

	return nil
}
