package main

import (
	"JD_Purchase/api"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
)

func main() {
	skuIDs := "730618,4080291:2"
	area := "18_1482_48938_52586"
	u, _ := url.Parse("socks5://127.0.0.1:8889")
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
	}
	purchaser, err := api.NewApi(client)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	result, err := purchaser.LoginByQRCode()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	log.Println(result)
	err = purchaser.BuyItemInStock(skuIDs, area, false, 5, 3, 5)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
}
