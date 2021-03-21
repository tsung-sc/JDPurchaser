package models

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Messenger struct {
	ScKey string
}

func NewMessenger(ScKey string) *Messenger {
	messenger := new(Messenger)
	messenger.ScKey = ScKey
	return messenger
}

func (m *Messenger) Send(client *http.Client, text string, desp string) error {
	nowTime := time.Now().Format(time.RFC3339)
	if desp != "" {
		desp = fmt.Sprintf("%s [%s]", desp, nowTime)
	} else {
		desp = fmt.Sprintf("[%s]", nowTime)
	}
	u := fmt.Sprintf("https://sc.ftqq.com/%s.send?text=%s&desp=%s", m.ScKey, url.QueryEscape(text), url.QueryEscape(desp))
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	errno := jsoniter.Get(data, "errno").ToInt64()
	if errno == 0 {
		log.Printf("Message sent successfully [text: %s, desp: %s]", text, desp)
	} else {
		log.Printf("Fail to send message, reason: %s", string(data))
	}
	return nil
}
