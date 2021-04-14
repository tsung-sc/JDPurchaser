package main

import (
	"JD_Purchase/api"
	"log"
	"net/http"
)

func main() {
	skuIDs := "730618,4080291:2"
	area := "18_1482_48938_52586"
	client := &http.Client{}
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
