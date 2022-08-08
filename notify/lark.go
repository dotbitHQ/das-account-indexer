package notify

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

type MsgContent struct {
	Tag      string `json:"tag"`
	UserId   string `json:"user_id,omitempty"`
	UnEscape bool   `json:"un_escape,omitempty"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
}
type MsgData struct {
	Email   string `json:"email"`
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string         `json:"title"`
				Content [][]MsgContent `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

func SendLarkTextNotify(url, title, text string) error {
	if url == "" || text == "" {
		return nil
	}
	var data MsgData
	data.Email = ""
	data.MsgType = "post"
	data.Content.Post.ZhCn.Title = title
	data.Content.Post.ZhCn.Content = [][]MsgContent{
		{
			MsgContent{
				Tag:      "text",
				UnEscape: false,
				Text:     text,
			},
		},
	}
	resp, _, errs := gorequest.New().Post(url).Timeout(time.Second * 10).SendStruct(&data).End()
	if len(errs) > 0 {
		return fmt.Errorf("errs:%v", errs)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%d", resp.StatusCode)
	}
	return nil
}
