package notify

import (
	"das-account-indexer/prometheus"
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

func SendLarkErrNotify(title, text string) {
	if title == "" || text == "" {
		return
	}
	prometheus.Tools.Metrics.ErrNotify().WithLabelValues(title, text).Inc()
}
