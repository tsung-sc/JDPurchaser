package models_test

import (
	"JD_Purchase/models"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestMessenger_Send(t *testing.T) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	client := http.Client{
		Transport: transport,
	}
	messenger := models.NewMessenger("")
	err := messenger.Send(&client, "jd-assistant 订单提交成功", "订单号：%s % order_id")
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}
