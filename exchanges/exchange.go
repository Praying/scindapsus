package exchanges

import (
	"scindapsus/config"
	"scindapsus/exchanges/okex"
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
	WatchOrderBook(symbol string, limit string, params interface{})
	WatchTicker(symbol string, limit string, params interface{})
	WatchTickers(symbols []string, params interface{})
	WatchOHLCV(symbol string, timeframe string, since string, limit string, params interface{})
	WatchTrades(symbol, since string, limit string, params interface{})
	//私有频道方法
	WatchBalance(params interface{})
	WatchOrders(symbol string, since string, limit string, params interface{})
	WatchCreateOrder(symbol, rtype, side string, amount, price float64, params interface{})
	WatchCancelOrder(id, symbol string, params interface{})
	WatchMyTrades(symbol, since, limit, params interface{})
	WatchPosition(symbol string)
	//Other
	Init()
}

func (Ok *OKExchange) WatchOrderBook(symbol string, limit string, params interface{}) {
	Ok.publicWS.WatchDepth([]string{symbol})
}

func (Ok *OKExchange) WatchTicker(symbol string, limit string, params interface{}) {
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

func (Ok *OKExchange) WatchOrders(symbol string, since string, limit string, params interface{}) {
	Ok.privateWS.WatchOrders(okex.MARGIN, symbol)
}

func (Ok *OKExchange) WatchCreateOrder(symbol, rtype, side string, amount, price float64, params interface{}) {
	Ok.privateWS.WatchCreateOrder(symbol, rtype, side, amount, price)
}

func (Ok *OKExchange) WatchCancelOrder(id, symbol string, params interface{}) {
	Ok.privateWS.WatchCancelOrder(id, symbol)
}

func (Ok *OKExchange) WatchMyTrades(symbol, since, limit, params interface{}) {
	//这是想看到交易量
	panic("implement me")
}

func (Ok *OKExchange) WatchPosition(symbol string) {
	Ok.privateWS.WatchPosition(okex.MARGIN, symbol)
}

type OKExchange struct {
	publicWS  *okex.OKExWSClient
	privateWS *okex.OKExWSClient
}

func NewOKExchange() *OKExchange {
	return &OKExchange{}
}

func (Ok *OKExchange) Init() {
	Ok.publicWS = okex.NewOKExWSClient(okex.TEST_PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
	Ok.publicWS.ConnectWS()
	Ok.privateWS = okex.NewOKExWSClient(okex.TEST_PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
	Ok.privateWS.ConnectWS()
	//登录
	apiConfig := config.GetConfigEngine().ReadConfig()
	Ok.privateWS.Login(apiConfig)
}
