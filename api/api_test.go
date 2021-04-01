package api_test

import (
	"JD_Purchase/api"
	"crypto/tls"
	"encoding/gob"
	"os"

	//"encoding/gob"
	"log"
	"net/http"
	"net/url"
	//"os"
	"testing"
)

func TestNewApi(t *testing.T) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	client := &http.Client{
		Transport: transport,
	}
	api, err := api.NewApi(client)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Log(api)
}

func TestApi_GetQRCode(t *testing.T) {
	test := new(api.Api)
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.GetQRCode()
	if err != nil {
		log.Println(err)
	}
	log.Println(ret)
}

func TestApi_GetQRCodeTicket(t *testing.T) {
	test := new(api.Api)
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.GetQRCodeTicket()
	if err != nil {
		log.Println(err)
	}
	log.Println(ret)
}

func TestApi_LoginByQRCode(t *testing.T) {
	proxyUrl, _ := url.Parse("socks5://127.0.0.1:8889")
	proxy := http.ProxyURL(proxyUrl)
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	client := &http.Client{
		Transport: transport,
	}
	test, err := api.NewApi(client)
	if err != nil {
		t.Logf("%+v", err)
		return
	}
	ret, err := test.LoginByQRCode()
	if err != nil {
		log.Printf("%+v", err)
	}
	log.Println(ret)
}

func TestApi_GetItemDetailPage(t *testing.T) {
	test := new(api.Api)
	test.Headers = make(http.Header, 1)
	test.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	//test.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.GetItemDetailPage("1179553")
	if err != nil {
		log.Printf("+%v", err)
	}
	log.Println(string(ret))
}

func TestApi_GetSigleItemStock(t *testing.T) {
	test := new(api.Api)
	test.ItemCat = make(map[string]string)
	test.ItemVenderIDs = make(map[string]string)
	test.Headers = make(http.Header, 1)
	test.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.GetSigleItemStock("10023774354580", "1", "2_2834_51993_0")
	if err != nil {
		log.Printf("+%v", err)
		return
	}
	log.Println(ret)
}

func TestApi_CancelSelectAllCartItem(t *testing.T) {
	test := new(api.Api)
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	test.Headers = make(http.Header, 1)
	test.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	//test.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.CancelSelectAllCartItem()
	if err != nil {
		log.Printf("%+v", err)
	}
	log.Println(ret)
}

func TestApi_AddItemToCart(t *testing.T) {
	var cookiee []*http.Cookie
	test := new(api.Api)

	f, _ := os.Open("cookies")
	bi := gob.NewDecoder(f)
	err := bi.Decode(&cookiee)
	if err != nil {
		log.Println(err)
	}
	log.Println(cookiee[1].Value)
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	//ret,err:=test.LoginByQRCode()
	//if err!=nil{
	//	log.Println(err)
	//}
	//log.Println(ret)
	////f,_:=os.Create("./cookies")
	////bi:=gob.NewEncoder(f)
	////u,_:=url.Parse("https://qr.m.jd.com")
	////err=bi.Encode(test.Sess.Cookies(u))
	////if err != nil {
	////	fmt.Println("编码错误", err.Error())
	////} else {
	////	fmt.Println("编码成功")
	////}
	//test.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	test.Headers = make(http.Header, 1)
	test.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	//test.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	test.AddItemToCart("100001324422", "1")
	//if err!=nil{
	//	log.Printf("%+v",err)
	//}
	//log.Println(ret)
}

func TestApi_GetMultiItemStockNew(t *testing.T) {
	itemMap := make(map[string]string)
	itemMap["730618"] = "1"
	itemMap["4080291"] = "2"
	test := new(api.Api)
	test.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	test.Headers = make(http.Header, 1)
	test.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	//test.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("socks5://127.0.0.1:8889")
	}
	transport := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	test.Client = &http.Client{
		Transport: transport,
	}
	ret, err := test.GetMultiItemStockNew(itemMap, "18_1482_48938_52586")
	if err != nil {
		log.Printf("%+v", err)
	}
	log.Println(ret)
}

func TestApi_(t *testing.T) {

}
