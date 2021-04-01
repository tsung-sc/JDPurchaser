package api

import (
	"JD_Purchase/config"
	"JD_Purchase/models"
	"JD_Purchase/utils"
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_TIMEOUT    = 10
	DEFAULT_USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36`
	useRandomUa        = ""
	QRCodeFile         = "./QRCode.png"
	retryTimes         = 85
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

type Api struct {
	Messenger        *models.Messenger
	Client           *http.Client
	UserAgent        string
	Headers          http.Header
	EID              string
	Fp               string
	TrackID          string
	RiskControl      string
	Timeout          time.Duration
	SendMessage      bool
	ItemCat          map[string]string
	ItemVenderIDs    map[string]string
	SeckillInitInfo  map[string]string
	SeckillOrderData map[string]string
	SeckillUrl       map[string]string
	Username         string
	Nickname         string
	IsLogin          bool
}

func NewApi(client *http.Client) (*Api, error) {
	var err error
	api := new(Api)
	api.Client = client
	if config.Get().IsRandomUserAgent {
		api.UserAgent = utils.GetRandomUserAgent()
	} else {
		api.UserAgent = DEFAULT_USER_AGENT
	}
	header := map[string][]string{
		"User-Agent": {api.UserAgent},
	}
	api.Headers = header
	api.EID = config.Get().EID
	api.Fp = config.Get().FP
	api.TrackID = config.Get().TrackID
	api.RiskControl = config.Get().RiskControl
	if api.EID == "" || api.Fp == "" || api.TrackID == "" || api.RiskControl == "" {
		return nil, xerrors.Errorf("请在 config.ini 中配置 eid, fp, track_id, risk_control 参数")
	}
	if config.Get().Timeout == 0 {
		api.Timeout = time.Duration(DEFAULT_TIMEOUT)
	} else {
		api.Timeout = time.Duration(config.Get().Timeout)
	}
	api.SendMessage = config.Get().EnableSendMessage
	if api.SendMessage {
		api.Messenger = models.NewMessenger(config.Get().Messenger.Sckey)
	}
	api.ItemCat = make(map[string]string)
	api.ItemVenderIDs = make(map[string]string)
	api.SeckillInitInfo = make(map[string]string)
	api.SeckillOrderData = make(map[string]string)
	api.SeckillUrl = make(map[string]string)
	api.Username = ""
	api.Nickname = "JD_Purchase"
	api.IsLogin = false
	api.Client.Jar, err = cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	err = api.loadCookies()
	if err != nil {
		if xerrors.Is(err, os.ErrNotExist) {
			log.Printf("未找到Cookies,将重新登陆！")
		} else {
			return nil, xerrors.Errorf("%w", err)
		}
	}
	return api, nil
}

func (a *Api) loadCookies() error {
	var cookies []*http.Cookie
	_, err := os.Stat("./cookies")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./cookies", os.ModePerm)
			if err != nil {
				return xerrors.Errorf("%w", err)
			}
		} else {
			return xerrors.Errorf("%w", err)
		}
	}
	cookiesFile := path.Join("./cookies", fmt.Sprintf("%s.json", a.Nickname))
	cookiesByte, err := ioutil.ReadFile(cookiesFile)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	u, err := url.Parse("https://www.jd.com")
	err = jsoniter.Unmarshal(cookiesByte, &cookies)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	a.Client.Jar.SetCookies(u, cookies)
	a.IsLogin, err = a.validateCookies()
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (a *Api) saveCookies(cookies []*http.Cookie) error {
	_, err := os.Stat("./cookies")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./cookies", os.ModePerm)
			if err != nil {
				return xerrors.Errorf("%w", err)
			}
		} else {
			return xerrors.Errorf("%w", err)
		}
	}
	cookiesFile := path.Join("./cookies", fmt.Sprintf("%s.json", a.Nickname))
	f, err := os.Create(cookiesFile)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	defer f.Close()
	for _, cookie := range cookies {
		cookie.Domain = ".jd.com"
		cookie.Path = "/"
	}
	cookiesByte, err := jsoniter.Marshal(cookies)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	_, err = f.Write(cookiesByte)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (a *Api) validateCookies() (bool, error) {
	u := "https://order.jd.com/center/list.action"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header = a.Headers

	clientPool := sync.Pool{New: func() interface{} {
		return *a.Client
	}}
	newClient := clientPool.Get()
	defer clientPool.Put(newClient)
	v, ok := newClient.(http.Client)
	if !ok {
		return false, xerrors.Errorf("%w", "not http client!")
	}
	v.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := v.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}

func (a *Api) LoginByQRCode() (bool, error) {
	var ticket string
	if a.IsLogin {
		log.Println("登陆成功")
		return true, nil
	}
	a.GetLoginPage()
	ok, err := a.GetQRCode()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("二维码下载失败")
	}
	for i := 0; i < retryTimes; i++ {
		ticket, err = a.GetQRCodeTicket()
		if err != nil {
			return false, xerrors.Errorf("%w", err)
		}
		if ticket != "" {
			break
		}
		time.Sleep(time.Second * 2)
		if i == 85 {
			return false, xerrors.Errorf("二维码过期，请重新获取扫描")
		}
	}
	ok, err = a.ValidateQRCodeTicket(ticket)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	if !ok {
		return false, xerrors.Errorf("二维码信息校验失败")
	}
	log.Println("二维码登陆成功")
	a.IsLogin = true
	a.Nickname, err = a.GetUserInfo(a.saveCookies)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	return true, nil
}

func (a *Api) GetLoginPage() {
	u := "https://passport.jd.com/new/login.aspx"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header = a.Headers
	resp, err := a.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Println(resp)
}

func (a *Api) GetQRCode() (bool, error) {
	u := "https://qr.m.jd.com/show?"
	args := url.Values{}
	args.Add("appid", "133")
	args.Add("size", "147")
	args.Add("t", strconv.FormatInt(time.Now().Unix()*1e3, 10))
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", "https://passport.jd.com/new/login.aspx")
	resp, err := a.Client.Do(req)
	if err != nil {
		log.Println("获取二维码失败")
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	err = utils.SaveImage(resp, QRCodeFile)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	log.Println("获取二维码成功，请打开京东APP扫描")
	return true, nil
}

func (a *Api) GetQRCodeTicket() (string, error) {
	var token string
	u := "https://qr.m.jd.com/check?"
	args := url.Values{}
	args.Add("appid", "133")
	args.Add("callback", fmt.Sprintf("jQuery%v", rand.Intn(9999999-1000000)+1000000))
	cookieUrl, err := url.Parse("https://qr.m.jd.com")
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	cookies := a.Client.Jar.Cookies(cookieUrl)
	for _, v := range cookies {
		if v.Name == "wlfstk_smdl" {
			token = v.Value
			break
		}
	}
	if token == "" {
		return "", xerrors.Errorf("获取token失败")
	}
	args.Add("token", token)
	args.Add("_", strconv.FormatInt(time.Now().Unix()*1e3, 10))
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", "https://passport.jd.com/new/login.aspx")
	resp, err := a.Client.Do(req)
	if err != nil {
		log.Println("获取二维码扫描结果异常")
		return "", xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	ret := new(models.QRCodeTicketBody)
	err = jsoniter.Unmarshal(data[14:len(data)-1], ret)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	if ret.Code != 200 {
		log.Printf("Code: %v,message: %s", ret.Code, ret.Msg)
		return "", nil
	}
	log.Println("已完成手机客户端确认！")
	return ret.Ticket, nil
}

func (a *Api) ValidateQRCodeTicket(ticket string) (bool, error) {
	u := "https://passport.jd.com/uc/qrCodeTicketValidation?"
	args := url.Values{}
	args.Add("t", ticket)
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", "https://passport.jd.com/uc/login?ltype=logout")
	resp, err := a.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	ret := new(models.ValidateQRCodeBody)
	err = jsoniter.Unmarshal(data, ret)
	if err != nil {
		return false, err
	}
	if ret.ReturnCode != 0 {
		return false, xerrors.Errorf("%w", ret.ReturnCode)
	}
	return true, nil
}

func (a *Api) GetUserInfo(cookiesHandle func([]*http.Cookie) error) (string, error) {
	u := "https://passport.jd.com/user/petName/getUserInfoForMiniJd.action?"
	args := url.Values{}
	args.Add("callback", fmt.Sprintf("jQuery%v", rand.Intn(9999999-1000000)+1000000))
	args.Add("_", fmt.Sprintf("%v", time.Now().Unix()*1e3))
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", "https://order.jd.com/center/list.action")
	resp, err := a.Client.Do(req)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	ret := new(models.UserInfo)
	err = jsoniter.Unmarshal(data[14:len(data)-1], ret)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	if ret.NickName == "" {
		return "jd", nil
	}
	err = cookiesHandle(req.Cookies())
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	return ret.NickName, nil
}

func (a *Api) BuyItemInStock(skuIDs string, area string, waitAll bool, stockInterval int, submitRetry int, submitInterval int) error {
	itemList := make([]string, 4)
	itemsMap := utils.ParseSkuID(skuIDs)
	areaID := utils.ParseAreaID(area)
	for k := range itemsMap {
		itemList = append(itemList, k)
	}
	if !waitAll {
		log.Printf("下单模式:%v 任一商品有货并且未下架均会尝试下单", itemList)
		for {
			for k, v := range itemsMap {
				ok, err := a.GetSigleItemStock(k, v, areaID)
				if err != nil {
					return xerrors.Errorf("%w", err)
				}
				if !ok {
					log.Printf("%s 不满足下单条件，%vs后进行下一次查询", k, stockInterval)
					continue
				}
				log.Printf("%s 满足下单条件，开始执行", k)
				_, _ = a.CancelSelectAllCartItem()
				a.AddOrChangeCartItem(map[string]string{}, k, v)
				ok, err = a.SubmitOrderWithRetry(submitRetry, submitInterval)
				if err != nil {
					return xerrors.Errorf("%w", err)
				}
				if ok {
					return nil
				}
				time.Sleep(time.Duration(stockInterval))
			}
		}
	} else {
		log.Printf("下单模式：%s 所有都商品同时有货并且未下架才会尝试下单", itemList)
		for {
			//todo
		}
	}
}

func (a *Api) GetMultiItemStockNew(itemMap map[string]string, area string) (bool, error) {
	u := "https://c0.3.cn/stocks?"
	keys := make([]string, 0)
	for k := range itemMap {
		keys = append(keys, k)
	}
	skuIds := strings.Join(keys, ",")
	args := url.Values{}
	args.Add("callback", fmt.Sprintf("jQuery%d", rand.Intn(9999999-1000000)+1000000))
	args.Add("type", "getstocks")
	args.Add("skuIds", skuIds)
	args.Add("area", area)
	args.Add("_", fmt.Sprintf("%d", time.Now().Unix()*1e3))
	req, err := http.NewRequest(http.MethodGet, u+args.Encode(), nil)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := a.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	stock := true
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	for _, key := range keys {
		skuState := jsoniter.Get(data[14:len(data)-1], key).Get("skuState").ToInt64()
		stockState := jsoniter.Get(data[14:len(data)-1], key).Get("StockState").ToInt()
		if skuState == 1 && (stockState >= 33 && stockState < 40) {
			continue
		} else {
			stock = false
			break
		}
	}
	return stock, nil
}

func (a *Api) GetSigleItemStock(skuID string, num string, area string) (bool, error) {
	var cat string
	var venderID string
	for k, v := range a.ItemCat {
		if k == skuID {
			cat = v
			break
		}
	}
	for k, v := range a.ItemVenderIDs {
		if k == skuID {
			venderID = v
			break
		}
	}
	if cat == "" {
		page, err := a.GetItemDetailPage(skuID)
		if err != nil {
			return false, xerrors.Errorf("%w", err)
		}
		reg := regexp.MustCompile(`cat: \[(.*?)\]`)
		cats := reg.FindStringSubmatch(page)
		cat = cats[1]
		a.ItemCat[skuID] = cat
		reg = regexp.MustCompile(`venderId:(\d*?),`)
		venderIDs := reg.FindStringSubmatch(page)
		venderID = venderIDs[1]
		a.ItemVenderIDs[skuID] = venderID
	}
	u := "https://c0.3.cn/stock?"
	args := url.Values{}
	args.Add("callback", fmt.Sprintf("jQuery%v", rand.Intn(9999999-1000000)+1000000))
	args.Add("buyNum", num)
	args.Add("skuId", skuID)
	args.Add("area", area)
	args.Add("_", fmt.Sprintf("%v", time.Now().Unix()*1e3))
	args.Add("ch", "1")
	args.Add("extraParam", "{\"originid\":\"1\"}")
	args.Add("cat", cat)
	args.Add("venderId", venderID)
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuID))
	a.Client.Timeout = time.Second * a.Timeout
	resp, err := a.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	skuState := jsoniter.Get(data[14:len(data)-1], "stock").Get("skuState").ToInt()
	stockState := jsoniter.Get(data[14:len(data)-1], "stock").Get("StockState").ToInt()
	if skuState != 1 || (stockState > 33 && stockState < 40) {
		return false, nil
	}
	return true, nil
}

func (a *Api) GetItemDetailPage(skuID string) (string, error) {
	u := fmt.Sprintf("https://item.jd.com/%s.html", skuID)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	req.Header = a.Headers
	resp, err := a.Client.Do(req)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	return string(data), nil
}

func (a *Api) CancelSelectAllCartItem() (bool, error) {
	u := "https://cart.jd.com/cancelAllItem.action"
	reqBody := new(models.CancelCartBody)
	reqBody.T = 0
	reqBody.OutSkus = ""
	reqBody.Random = rand.Intn(1)
	reqBodyByte, err := jsoniter.Marshal(reqBody)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(reqBodyByte))
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := a.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Panicf("Status: %d, Url: %s", resp.StatusCode, u)
		return false, nil
	}
	return true, nil
}

func (a *Api) GetCartDetail() (map[string]string, error) {
	catDetail := make(map[string]string)
	u := "https://cart.jd.com/cart.action"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	doc.Find("div[class=item-item]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
	return catDetail, xerrors.Errorf("%w", err)
}
func (a *Api) AddOrChangeCartItem(cart map[string]string, skuId string, count string) {
	for k, v := range cart {
		if skuId == v {
			log.Panicf("%s 已在购物车中，调整数量为 %s", skuId, count)
			cartItem := v
			log.Println(k, cartItem)
			//todo:
		}
	}
	log.Printf("%s 不在购物车中，开始加入购物车,数量 %s", skuId, count)
	a.AddItemToCart(skuId, count)
}

func (a *Api) AddItemToCart(skuId string, count string) {
	var result bool
	u := "https://cart.jd.com/gate.action?"
	args := url.Values{}
	args.Add("pid", skuId)
	args.Add("pcount", count)
	args.Add("ptype", "1")
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		//return xerrors.Errorf("%w",err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		//return nil, xerrors.Errorf("%w", err)
	}
	doc.Find("h3[class=ftx-02]").Each(func(i int, s *goquery.Selection) {
		result = s.Text() == "商品已成功加入购物车！"
	})
	if result {
		log.Printf("%s x %s已成功加入购物车", skuId, count)
	} else {
		log.Printf("%s 添加到购物车失败", skuId)
	}
}

func (a *Api) SubmitOrderWithRetry(submitRetry, interval int) (bool, error) {
	for i := 1; i < submitRetry+1; i++ {
		log.Printf("第[%d/%d]次尝试提交订单", i, submitRetry)
		a.GetCheckoutPageDetail()
		result, err := a.SubmitOrder()
		if err != nil {
			return false, xerrors.Errorf("%w", err)
		}
		if result {
			log.Printf("第%d次提交订单成功", i)
			return true, nil
		} else {
			if i < retryTimes {
				log.Printf("第%d次提交失败，%ds后重试", i, interval)
				time.Sleep(time.Second * time.Duration(interval))
			}
		}
	}
	log.Printf("重试提交%d次结束", retryTimes)
	return false, nil
}

func (a *Api) GetCheckoutPageDetail() (map[string]string, error) {
	orderDetail := make(map[string]string)
	u := "http://trade.jd.com/shopping/order/getOrderInfo.action?"
	args := url.Values{}
	args.Add("rid", fmt.Sprintf("%d", time.Now().Unix()*1e3))
	u = u + args.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("获取订单结算页信息失败")
		return nil, xerrors.Errorf("%w", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	orderDetail["address"] = strings.Trim(doc.Find("span[id=sendAddr]").Text(), "寄送至： ")
	orderDetail["receiver"] = strings.Trim(doc.Find("span[id=sendMobile]").Text(), "收货人：")
	orderDetail["total_price"] = strings.Trim(doc.Find("span[id=sumPayPriceId]").Text(), "￥")
	//orderDetail["items"]=[]string{}
	log.Printf("下单信息:%v", orderDetail)
	return orderDetail, nil
}

func (a *Api) SubmitOrder() (bool, error) {
	var orderId int64
	u := "https://trade.jd.com/shopping/order/submitOrder.action"
	args := url.Values{}
	args.Add("overseaPurchaseCookies", "")
	args.Add("vendorRemarks", "[]")
	args.Add("submitOrderParam.sopNotPutInvoice", "false")
	args.Add("submitOrderParam.trackID", "TestTrackId")
	args.Add("submitOrderParam.ignorePriceChange", "0")
	args.Add("submitOrderParam.btSupport", "0")
	args.Add("riskControl", a.RiskControl)
	args.Add("submitOrderParam.isBestCoupon", "1")
	args.Add("submitOrderParam.jxj", "1")
	args.Add("submitOrderParam.trackId", a.TrackID)
	args.Add("submitOrderParam.eid", a.EID)
	args.Add("submitOrderParam.fp", a.Fp)
	args.Add("submitOrderParam.needCheck", "1")

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(args.Encode()))
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("Host", "trade.jd.com")
	req.Header.Add("Referer", "http://trade.jd.com/shopping/order/getOrderInfo.action")
	resp, err := a.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, xerrors.Errorf("%w", err)
	}
	if jsoniter.Get(data, "success").ToBool() {
		orderId = jsoniter.Get(data, "orderId").ToInt64()
		log.Printf("订单提交成功！订单号:%d", orderId)
		if a.SendMessage {
			//todo:

		}
		return true, nil
	} else {
		message := jsoniter.Get(data, "message").ToString()
		resultCode := jsoniter.Get(data, "resultCode").ToInt64()
		switch resultCode {
		case 0:
			a.SaveInvoice()
			message = message + "(下单商品可能为第三方商品，将切换为普通发票进行尝试)"
		case 60077:
			message = message + "(可能是购物车为空 或 未勾选购物车中商品)"
		case 60123:
			message = message + "(需要在config.ini文件中配置支付密码)"
		default:
			log.Println(resultCode)
		}
		log.Printf("订单提交失败, 错误码：%d, 返回信息：%s", resultCode, message)
		return false, nil
	}
}

func (a *Api) SaveInvoice() {
	u := "https://trade.jd.com/shopping/dynamic/invoice/saveInvoice.action"
	args := url.Values{}
	args.Add("invoiceParam.selectedInvoiceType", "1")
	args.Add("invoiceParam.companyName", "个人")
	args.Add("invoiceParam.invoicePutType", "0")
	args.Add("invoiceParam.selectInvoiceTitle", "4")
	args.Add("invoiceParam.selectBookInvoiceContent", "")
	args.Add("invoiceParam.selectNormalInvoiceContent", "1")
	args.Add("invoiceParam.vatCompanyName", "")
	args.Add("invoiceParam.code", "")
	args.Add("invoiceParam.regAddr", "")
	args.Add("invoiceParam.regPhone", "")
	args.Add("invoiceParam.regBank", "")
	args.Add("invoiceParam.regBankAccount", "1")
	args.Add("invoiceParam.hasCommon", "true")
	args.Add("invoiceParam.hasBook", "false")
	args.Add("invoiceParam.consigneeName", "")
	args.Add("invoiceParam.consigneePhone", "")
	args.Add("invoiceParam.consigneeAddress", "")
	args.Add("invoiceParam.consigneeProvince", "请选择：")
	args.Add("invoiceParam.consigneeProvinceId", "NaN")
	args.Add("invoiceParam.consigneeCity", "请选择")
	args.Add("invoiceParam.consigneeCityId", "NaN")
	args.Add("invoiceParam.consigneeCounty", "请选择")
	args.Add("invoiceParam.consigneeCountyId", "NaN")
	args.Add("invoiceParam.consigneeTown", "请选择")
	args.Add("invoiceParam.consigneeTownId", "0")
	args.Add("invoiceParam.sendSeparate", "false")
	args.Add("invoiceParam.usualInvoiceId", "")
	args.Add("invoiceParam.selectElectroTitle", "4")
	args.Add("invoiceParam.electroCompanyName", "undefined")
	args.Add("invoiceParam.electroInvoiceEmail", "")
	args.Add("invoiceParam.electroInvoicePhone", "")
	args.Add("invokeInvoiceBasicService", "true")
	args.Add("invoice_ceshi1", "")
	args.Add("invoiceParam.showInvoiceSeparate", "false")
	args.Add("invoiceParam.invoiceSeparateSwitch", "1")
	args.Add("invoiceParam.invoiceCode", "")
	args.Add("invoiceParam.saveInvoiceFlag", "1")
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(args.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Referer", "https://trade.jd.com/shopping/dynamic/invoice/saveInvoice.action")
	resp, err := a.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
