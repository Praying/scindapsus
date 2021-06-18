package okex

//推送数据参数
type TickerResp struct {
	Arg struct {
		Channel string `json:"channel"`
		Instid  string `json:"instId"`
	} `json:"arg"`
	Data []struct {
		Insttype  string `json:"instType"`
		Instid    string `json:"instId"`
		Last      string `json:"last"`
		Lastsz    string `json:"lastSz"`
		Askpx     string `json:"askPx"`
		Asksz     string `json:"askSz"`
		Bidpx     string `json:"bidPx"`
		Bidsz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		Volccy24H string `json:"volCcy24h"` //Ccy是以币为单位
		Vol24H    string `json:"vol24h"`
		Sodutc0   string `json:"sodUtc0"`
		Sodutc8   string `json:"sodUtc8"`
		Ts        string `json:"ts"`
	} `json:"data"`
}

//深度数据
type BookResp struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Action string `json:"action"`
	Data   []struct {
		Asks     [][]string `json:"asks"`
		Bids     [][]string `json:"bids"`
		Ts       string     `json:"ts"`
		Checksum int        `json:"checksum"`
	} `json:"data"`
}

type SubParam struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"args"`
}

//私有频道，登录
type LoginParam struct {
	Op   string `json:"op"`
	Args []struct {
		APIKey     string `json:"apiKey"`
		Passphrase string `json:"passphrase"`
		Timestamp  string `json:"timestamp"`
		Sign       string `json:"sign"`
	} `json:"args"`
}

//私有频道，订阅持仓 参数
type PositionParam struct {
	Op   string `json:"op"`
	Args []struct {
		Channel  string `json:"channel"`  //必填
		InstType string `json:"instType"` //必填
		Uly      string `json:"uly"`      //	合约标的指数
		InstID   string `json:"instId"`   //产品ID
	} `json:"args"`
}

type PositionResp struct {
	Arg struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
	} `json:"arg"`
	Data []struct {
		Adl      string `json:"adl"`
		AvailPos string `json:"availPos"`
		AvgPx    string `json:"avgPx"`
		CTime    string `json:"cTime"`
		Ccy      string `json:"ccy"`
		DeltaBS  string `json:"deltaBS"`
		DeltaPA  string `json:"deltaPA"`
		GammaBS  string `json:"gammaBS"`
		GammaPA  string `json:"gammaPA"`
		Imr      string `json:"imr"`
		InstID   string `json:"instId"`
		InstType string `json:"instType"`
		Interest string `json:"interest"`
		Last     string `json:"last"`
		Lever    string `json:"lever"`
		Liab     string `json:"liab"`
		LiabCcy  string `json:"liabCcy"`
		LiqPx    string `json:"liqPx"`
		Margin   string `json:"margin"`
		MgnMode  string `json:"mgnMode"`
		MgnRatio string `json:"mgnRatio"`
		Mmr      string `json:"mmr"`
		OptVal   string `json:"optVal"`
		PTime    string `json:"pTime"`
		Pos      string `json:"pos"`
		PosCcy   string `json:"posCcy"`
		PosID    string `json:"posId"`
		PosSide  string `json:"posSide"`
		ThetaBS  string `json:"thetaBS"`
		ThetaPA  string `json:"thetaPA"`
		TradeID  string `json:"tradeId"`
		UTime    string `json:"uTime"`
		Upl      string `json:"upl"`
		UplRatio string `json:"uplRatio"`
		VegaBS   string `json:"vegaBS"`
		VegaPA   string `json:"vegaPA"`
	} `json:"data"`
}

//下单的订单请求
type OrderParam struct {
	ID   string `json:"id"`
	Op   string `json:"op"`
	Args []struct {
		ClOrderID string `json:"clOrdId"`
		Side      string `json:"side"`
		InstID    string `json:"instId"`
		TdMode    string `json:"tdMode"`
		OrdType   string `json:"ordType"`
		Sz        string `json:"sz"` //限价单时，表示数量
		Px        string `json:"px"` //限价单，表示委托价格
	} `json:"args"`
}

//余额和持仓
type BalAndPosResp struct {
	Arg struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Data []struct {
		PTime     string `json:"pTime"`
		EventType string `json:"eventType"`
		BalData   []struct {
			Ccy     string `json:"ccy"`
			CashBal string `json:"cashBal"`
			UTime   string `json:"uTime"`
		} `json:"balData"`
		PosData []struct {
			PosID    string `json:"posId"`
			TradeID  string `json:"tradeId"`
			InstID   string `json:"instId"`
			InstType string `json:"instType"`
			MgnMode  string `json:"mgnMode"`
			PosSide  string `json:"posSide"`
			Pos      string `json:"pos"`
			Ccy      string `json:"ccy"`
			PosCcy   string `json:"posCcy"`
			AvgPx    string `json:"avgPx"`
			UTime    string `json:"uTime"`
		} `json:"posData"`
	} `json:"data"`
}

//WS下单操作的响应
type OrderResp struct {
	Code string `json:"code"`
	Data []struct {
		ClOrdID string `json:"clOrdId"`
		OrdID   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
		Tag     string `json:"tag"`
	} `json:"data"`
	ID  string `json:"id"`
	Msg string `json:"msg"`
	Op  string `json:"op"`
}
