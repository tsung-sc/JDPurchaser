package models

type Stock struct {
	Stock struct {
		FreshEdi  interface{} `json:"freshEdi"`
		Ext       string      `json:"Ext"`
		RealSkuID int         `json:"realSkuId"`
		Pr        struct {
			PromiseResult string `json:"promiseResult"`
			ResultCode    int    `json:"resultCode"`
		} `json:"pr"`
		PromiseResult     string        `json:"promiseResult"`
		NationallySetWare string        `json:"nationallySetWare"`
		SidDely           string        `json:"sidDely"`
		Channel           int           `json:"channel"`
		ServiceInfo       string        `json:"serviceInfo"`
		Rid               interface{}   `json:"rid"`
		WeightValue       string        `json:"weightValue"`
		IsSopUseSelfStock string        `json:"isSopUseSelfStock"`
		IsPurchase        bool          `json:"IsPurchase"`
		AreaLevel         int           `json:"areaLevel"`
		Cla               []interface{} `json:"cla"`
		Eb                string        `json:"eb"`
		Ec                string        `json:"ec"`
		SkuID             int           `json:"skuId"`
		Dc                []interface{} `json:"Dc"`
		YjhxMap           []interface{} `json:"yjhxMap"`
		IsWalMar          bool          `json:"isWalMar"`
		Area              struct {
			TownName     string `json:"townName"`
			CityName     string `json:"cityName"`
			Success      bool   `json:"success"`
			ProvinceName string `json:"provinceName"`
			CountyName   string `json:"countyName"`
		} `json:"area"`
		Ab          string `json:"ab"`
		Ac          string `json:"ac"`
		Ad          string `json:"ad"`
		Ae          string `json:"ae"`
		SkuState    int    `json:"skuState"`
		PopType     int    `json:"PopType"`
		Af          string `json:"af"`
		PromiseMark string `json:"promiseMark"`
		Ag          string `json:"ag"`
		IsSam       bool   `json:"isSam"`
		Ir          []struct {
			ResultCode      int    `json:"resultCode"`
			ShowName        string `json:"showName"`
			HelpLink        string `json:"helpLink"`
			IconTip         string `json:"iconTip"`
			PicURL          string `json:"picUrl"`
			IconCode        string `json:"iconCode"`
			IconType        int    `json:"iconType"`
			IconSrc         string `json:"iconSrc"`
			IconServiceType int    `json:"iconServiceType"`
		} `json:"ir"`
		Vd      interface{} `json:"vd"`
		Rfg     int         `json:"rfg"`
		Dti     interface{} `json:"Dti"`
		JdPrice struct {
			P  string `json:"p"`
			Op string `json:"op"`
			ID string `json:"id"`
			M  string `json:"m"`
		} `json:"jdPrice"`
		Rn             int           `json:"rn"`
		Support        []interface{} `json:"support"`
		VenderType     interface{}   `json:"venderType"`
		PlusFlagInfo   string        `json:"PlusFlagInfo"`
		Code           int           `json:"code"`
		IsPlus         bool          `json:"isPlus"`
		AfsCode        int           `json:"afsCode"`
		StockDesc      string        `json:"stockDesc"`
		Sid            string        `json:"sid"`
		IsJDexpress    string        `json:"isJDexpress"`
		DcID           string        `json:"dcId"`
		Sr             interface{}   `json:"sr"`
		StockState     int           `json:"StockState"`
		StockStateName string        `json:"StockStateName"`
		DcashDesc      string        `json:"dcashDesc"`
		M              string        `json:"m"`
		SelfD          struct {
			Vid       int         `json:"vid"`
			Df        interface{} `json:"df"`
			Cg        string      `json:"cg"`
			ColType   int         `json:"colType"`
			Deliver   string      `json:"deliver"`
			ID        int         `json:"id"`
			Type      int         `json:"type"`
			Vender    string      `json:"vender"`
			Linkphone string      `json:"linkphone"`
			URL       string      `json:"url"`
			Po        string      `json:"po"`
		} `json:"self_D"`
		ArrivalDate string `json:"ArrivalDate"`
		V           string `json:"v"`
		PromiseYX   struct {
			ResultCode      int    `json:"resultCode"`
			ShowName        string `json:"showName"`
			HelpLink        string `json:"helpLink"`
			IconTip         string `json:"iconTip"`
			PicURL          string `json:"picUrl"`
			IconCode        string `json:"iconCode"`
			IconType        int    `json:"iconType"`
			IconSrc         string `json:"iconSrc"`
			IconServiceType int    `json:"iconServiceType"`
		} `json:"promiseYX"`
	} `json:"stock"`
	ChoseSuit []interface{} `json:"choseSuit"`
}
