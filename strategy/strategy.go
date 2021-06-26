package strategy

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"scindapsus/exchanges/okex"
	"time"
)

type AbstractStrategy interface {
	Symbol() string
	OnTicker(tickerData bd.TickerData)
	Init()
	Name() string
}

//现货网格马丁策略
type SpotGridMartinStrategy struct {
	pubWS  *okex.OKExWSClient
	privWS *okex.OKExWSClient

	//余额，USDT
	balance float64
	//持仓，ETH
	position float64

	ssymbol string

	//策略名称
	StName string

	Inited bool
}

func NewSpotGridMartinStrategy(symbol string) *SpotGridMartinStrategy {
	name := symbol + time.ANSIC
	return &SpotGridMartinStrategy{
		ssymbol: symbol,
		StName:  name,
		Inited:  false,
	}
}

func (strategy *SpotGridMartinStrategy) Name() string {
	return strategy.StName
}

func (strategy *SpotGridMartinStrategy) Init() {
	strategy.pubWS = okex.NewOKExWSClient(okex.TEST_PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
	strategy.pubWS.ConnectWS()
	strategy.privWS = okex.NewOKExWSClient(okex.TEST_PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
	strategy.privWS.ConnectWS()
	apiConfig := &okex.APIConfig{
		HttpClient:    nil,
		Endpoint:      "",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
		ClientId:      "",
		Lever:         0,
	}
	strategy.privWS.Login(apiConfig)
	strategy.pubWS.SubscribeTicker([]string{"ETH-USDT"})
	strategy.Inited = true
	log.Infof("[%s] inited", strategy.StName)
}

func (strategy *SpotGridMartinStrategy) Symbol() string {
	return strategy.ssymbol
}

func (strategy *SpotGridMartinStrategy) OnTicker(tickerData bd.TickerData) {
	log.Info("[SpotGridMartinStrategy]")
	//
}

func (strategy *SpotGridMartinStrategy) OnPosition(data bd.TickerData) {

}
