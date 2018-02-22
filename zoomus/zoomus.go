package zoomus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ZoomClient struct {
	WebhookURL *url.URL
	HTTPClient *http.Client
	Header     map[string]string
}

type Message struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
	Action  string `json:"action"`
}

func NewZoomClient(webhook, token string) (*ZoomClient, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("token is missing")
	}

	p, err := url.Parse(webhook)
	if err != nil {
		return nil, err
	}
	zc := &ZoomClient{
		WebhookURL: p,
		HTTPClient: &http.Client{},
		Header: map[string]string{
			"Content-Type": "application/json",
			"X-Zoom-Token": token,
		},
	}

	return zc, nil
}

func makeJSONMassage(msg *Message) ([]byte, error) {
	format := "<p><label>%s</label></p>"
	msg.Title = fmt.Sprintf(format, msg.Title)
	msg.Summary = fmt.Sprintf(format, msg.Summary)
	msg.Body = fmt.Sprintf(format, msg.Body)
	msg.Action = fmt.Sprintf(map[string]string{
		"send": "<p><button onclick=\"sendMsg('1', %s)\">send</button></p>",
		"copy": "<p><button onclick=\"copyMsg('1', %s)\">copy</button></p>",
	}[msg.Action], msg.Action)
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return msgJSON, nil
}

func (zc *ZoomClient) SendMessage(msg *Message) error {
	msgJSON, err := makeJSONMassage(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", zc.WebhookURL.String(), bytes.NewBuffer(msgJSON))
	for k, v := range zc.Header {
		req.Header.Set(k, v)
	}
	res, err := zc.HTTPClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed: status: %s: url=%s", res.Status, zc.WebhookURL.String())
	}
	return nil
}
