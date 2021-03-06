package exchanges

import (
	"scindapsus/config"
	"scindapsus/exchanges/okex"
	"scindapsus/util"
	"time"
)

/*

copy from ccxt
watchOrderBook (symbol, limit, params)
watchTicker (symbol, limit, params)
watchTickers (symbols, params) wip
watchOHLCV (symbol, timeframe, since, limit, params)
watchTrades (symbol, since, limit, params)

watchBalance (params)
watchOrders (symbol, since, limit, params) wip
watchCreateOrder (symbol, type, side, amount, price, params) wip
watchCancelOrder (id, symbol, params) wip
watchMyTrades (symbol, since, limit, params) wip
*/
type Exchange interface {
	//公共频道方法
	WatchInstruments(insTypes []string)
	WatchOrderBook(symbol string, limit string, params interface{})
	WatchTicker(symbol string)
	WatchTickers(symbols []string, params interface{})
	WatchOHLCV(symbol string, timeframe string, since string, limit string, params interface{})
	WatchTrades(symbol, since string, limit string, params interface{})
	//订阅资金费率
	WatchFundingRate(symbol string)
	//私有频道方法
	WatchBalance(params interface{})
	WatchOrders(symbol string, instType string)
	WatchCreateOrder(symbol, rtype, side string, amount, price float64, tradeMode string)
	WatchCancelOrder(id, symbol string, params interface{})
	WatchMyTrades(symbol, since, limit, params interface{})
	WatchPosition(symbol string)
	//Other
	Init()
}

func (Ok *OKExchange) WatchFundingRate(symbol string) {
	Ok.publicWS.WatchFundingRate(symbol)
}

func (Ok *OKExchange) WatchOrderBook(symbol string, limit string, params interface{}) {
	Ok.publicWS.WatchDepth([]string{symbol})
}

func (Ok *OKExchange) WatchTicker(symbol string) {
	Ok.publicWS.WatchTicker([]string{symbol})
}

func (Ok *OKExchange) WatchTickers(symbols []string, params interface{}) {
	Ok.publicWS.WatchTicker(symbols)
}

func (Ok *OKExchange) WatchOHLCV(symbol string, timeframe string, since string, limit string, params interface{}) {
	panic("implement me")
}

//公共频道的方法
func (Ok *OKExchange) WatchTrades(symbol, since string, limit string, params interface{}) {
	Ok.publicWS.WatchTrades(symbol)
}

func (Ok *OKExchange) WatchBalance(params interface{}) {
	Ok.privateWS.WatchBalAndPos()
}

func (Ok *OKExchange) WatchOrders(symbol string, instType string) {
	if okex.CheckInstType(instType) {
		Ok.privateWS.WatchOrders(instType, symbol)
	} else {
		Ok.privateWS.WatchOrders(okex.INST_ANY, symbol)
	}

}

func (Ok *OKExchange) WatchCreateOrder(symbol, rtype, side string, amount, price float64, tradeMode string) {
	clOrdId := util.GenerateClOrdId(Ok.ConnectTime, Ok.OrderCount)
	Ok.OrderCount = Ok.OrderCount + 1
	Ok.privateWS.WatchCreateOrder(symbol, rtype, side, amount, price, clOrdId, tradeMode)
}

func (Ok *OKExchange) WatchCancelOrder(id, symbol string, params interface{}) {
	Ok.privateWS.WatchCancelOrder(id, symbol)
}

func (Ok *OKExchange) WatchMyTrades(symbol, since, limit, params interface{}) {
	//Ok.privateWS.watch
}

func (Ok *OKExchange) WatchPosition(symbol string) {
	Ok.privateWS.WatchPosition(okex.INST_ANY, symbol)
}

type OKExchange struct {
	publicWS  *okex.OKExWSClient
	privateWS *okex.OKExWSClient
	//生成订单的方式为 建立连接的时间+OrderCount
	OrderCount  int64
	ConnectTime int64
	//是否是实盘
	Real bool
}

func NewOKExchange() *OKExchange {
	return &OKExchange{
		OrderCount: 0,
		Real:       false,
	}
}

func (Ok *OKExchange) Init() {
	if Ok.Real {
		Ok.publicWS = okex.NewPublicWSClient(okex.PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
	} else {
		Ok.publicWS = okex.NewPublicWSClient(okex.TEST_PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
	}
	Ok.publicWS.ConnectWS()
	if Ok.Real {
		Ok.privateWS = okex.NewPrivateWSClient(okex.PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
	} else {
		Ok.privateWS = okex.NewPrivateWSClient(okex.TEST_PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
	}

	Ok.privateWS.ConnectWS()
	//登录
	apiConfig := config.GetConfigEngine().ReadConfig()
	Ok.privateWS.Login(apiConfig)
	Ok.ConnectTime = time.Now().Unix()
	time.Sleep(4 * time.Second)
}

func (Ok *OKExchange) WatchInstruments(insTypes []string) {
	Ok.publicWS.WatchInstrument(insTypes)
}
