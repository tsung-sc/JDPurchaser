package api

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//获取商品的抢购链接
//点击"抢购"按钮后，会有两次302跳转，最后到达订单结算页面
//这里返回第一次跳转后的页面url，作为商品的抢购链接
func (a *Api) GetSeckillUrl(skuID string) (string, error) {
	u := "https://itemko.jd.com/itemShowBtn?"
	args := url.Values{}
	args.Add("callback", fmt.Sprintf("jQuery%d", rand.Intn(9999999-1000000)+1000000))
	args.Add("skuId", skuID)
	args.Add("from", "pc")
	args.Add("_", fmt.Sprintf("%d", time.Now().Unix()*1e3))
	req, err := http.NewRequest(http.MethodGet, u+args.Encode(), nil)
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("Host", "itemko.jd.com")
	req.Header.Add("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuID))
	for {
		resp, err := a.Client.Do(req)
		if err != nil {
			return "", xerrors.Errorf("%w", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return "", xerrors.Errorf("%w", err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", xerrors.Errorf("%w", err)
		}
		u = jsoniter.Get(data, "url").ToString()
		if u != "" {
			routerUrl := fmt.Sprintf("https:%s", u)
			seckillUrl := strings.ReplaceAll(routerUrl, "divide", "marathon")
			seckillUrl = strings.ReplaceAll(seckillUrl, "user_routing", "captcha,html")
			log.Printf("抢购链接获取成功：%s", seckillUrl)
			return seckillUrl, nil
		} else {
			log.Printf("抢购链接获取失败，%s不是抢购商品或抢购页面暂未刷新,%v秒后重试", skuID, 0.5)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

//访问商品的抢购链接（用于设置cookie等）
func (a *Api) RequestSeckillUrl(skuID string) error {
	_, exist := a.SeckillUrl[skuID]
	if !exist {
		seckillUrl, err := a.GetSeckillUrl(skuID)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		a.SeckillUrl[skuID] = seckillUrl
	}
	req, err := http.NewRequest(http.MethodGet, a.SeckillUrl[skuID], nil)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("Host", "marathon.jd.com")
	req.Header.Add("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuID))

	a.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		a.Client.CheckRedirect = nil
	}()
	resp, err := a.Client.Do(req)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	return nil
}

//访问抢购订单结算页面
func (a *Api) RequestSeckillCheckoutPage(skuID string, num int) error {
	if num == 0 {
		num = 1
	}
	u := "https://marathon.jd.com/seckill/seckill.action?"
	args := url.Values{}
	args.Add("skuId", skuID)
	args.Add("num", fmt.Sprintf("%d", num))
	args.Add("rid", fmt.Sprintf("%d", time.Now().Unix()))
	req, err := http.NewRequest(http.MethodGet, u+args.Encode(), nil)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("Host", "marathon.jd.com")
	req.Header.Add("Referer", fmt.Sprintf("https://item.jd.com/%s.html", skuID))
	resp, err := a.Client.Do(req)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	return nil
}

//获取秒杀初始化信息（包括：地址，发票，token）
func (a *Api) GetSeckillInitInfo(skuID string, num int) (data []byte, err error) {
	if num == 0 {
		num = 1
	}
	u := "https://marathon.jd.com/seckillnew/orderService/pc/init.action"
	args := url.Values{}
	args.Add("skuId", skuID)
	args.Add("num", fmt.Sprintf("%d", num))
	args.Add("isModifyAddress", "false")
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(args.Encode()))
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("Host", "marathon.jd.com")
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return data, nil
}

//生成提交抢购订单所需的请求体参数
func (a *Api) GenSeckillOrderData(skuID string, num int) (string, error) {
	if num == 0 {
		num = 1
	}
	args := url.Values{}
	_, exist := a.SeckillInitInfo[skuID]
	if !exist {
		seckillInfo, err := a.GetSeckillInitInfo(skuID, 1)
		if err != nil {
			return "", xerrors.Errorf("%w", err)
		}
		a.SeckillInitInfo[skuID] = seckillInfo
	}
	initInfo := a.SeckillInitInfo[skuID]
	defaultAddress := jsoniter.Get(initInfo, "addressList", 0)
	invoiceInfoKeys := jsoniter.Get(initInfo, "invoiceInfo").Keys()
	if invoiceInfoKeys == nil {
		args.Add("invoiceTitle", "-1")
		args.Add("invoiceContent", "1")
		args.Add("invoicePhone", "")
		args.Add("invoicePhoneKey", "")
		args.Add("invoice", "false")
	} else {
		args.Add("invoiceTitle", jsoniter.Get(initInfo, "invoiceInfo").Get("invoiceTitle").ToString())
		args.Add("invoiceContent", jsoniter.Get(initInfo, "invoiceInfo").Get("invoiceContentType").ToString())
		args.Add("invoicePhone", jsoniter.Get(initInfo, "invoiceInfo").Get("invoicePhone").ToString())
		args.Add("invoicePhoneKey", jsoniter.Get(initInfo, "invoiceInfo").Get("invoicePhone").ToString())
		args.Add("invoice", "true")
	}
	token := jsoniter.Get(initInfo, "token").ToString()
	args.Add("skuId", skuID)
	args.Add("num", fmt.Sprintf("%d", num))
	args.Add("addressId", defaultAddress.Get("id").ToString())
	var yuShou string
	if jsoniter.Get(initInfo, "seckillSkuVO").Get("extMap").Get("YuShou").ToString() != "0" {
		yuShou = "true"
	} else {
		yuShou = "false"
	}
	args.Add("yuShou", yuShou)
	args.Add("isModifyAddress", "false")
	args.Add("name", defaultAddress.Get("name").ToString())
	args.Add("provinceId", defaultAddress.Get("provinceId").ToString())
	args.Add("cityId", defaultAddress.Get("cityId").ToString())
	args.Add("countyId", defaultAddress.Get("countyId").ToString())
	args.Add("townId", defaultAddress.Get("townId").ToString())
	args.Add("addressDetail", defaultAddress.Get("addressDetail").ToString())
	args.Add("mobile", defaultAddress.Get("mobile").ToString())
	args.Add("mobileKey", defaultAddress.Get("mobileKey").ToString())
	args.Add("email", defaultAddress.Get("email").ToString())
	args.Add("postCode", "")
	args.Add("invoiceCompanyName", "")
	args.Add("invoiceTaxpayerNO", "")
	args.Add("invoiceEmail", "")
	args.Add("password", a.PaymentPwd)
	args.Add("codTimeType", "3")
	args.Add("paymentType", "4")
	args.Add("areaCode", "")
	args.Add("overseas", "")
	args.Add("phone", "")
	args.Add("eid", a.EID)
	args.Add("fp", a.Fp)
	args.Add("token", token)
	args.Add("pru", "")
	return args.Encode(), nil
}
