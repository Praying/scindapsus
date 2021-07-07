package strategy

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"scindapsus/config"
	"scindapsus/exchanges/okex"
	"time"
)

type AbstractStrategy interface {
	Symbol() string
	OnTicker(tickerData bd.TickerData)
	OnPosition(positionData bd.PositionData)
	OnBalAndPos(balAndPosData bd.BalAndPosData)
	Init(apiConfig *config.APIConfig)
	Name() string
	IsInited() bool
	OnStart()
	IsTrading() bool
	OnStop()
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

	Trading bool
}

func NewSpotGridMartinStrategy(symbol string) *SpotGridMartinStrategy {
	name := symbol + time.ANSIC
	return &SpotGridMartinStrategy{
		ssymbol: symbol,
		StName:  name,
		Inited:  false,
	}
}

func (strategy *SpotGridMartinStrategy) OnBalAndPos(balAndPosData bd.BalAndPosData) {
	if bal, ok := balAndPosData.BalMap["USDT"]; ok {
		strategy.balance = bal
		log.Infof("[%s] update balance: %v", strategy.Name(), strategy.balance)
	}

}

func (strategy *SpotGridMartinStrategy) OnPosition(positionData bd.PositionData) {
	strategy.position = positionData.Pos
	log.Infof("[%s] update position:", strategy.Name(), positionData.Pos)
}

func (strategy *SpotGridMartinStrategy) OnStart() {
	strategy.Trading = true
}

func (strategy *SpotGridMartinStrategy) OnStop() {
	strategy.Trading = false
}

func (strategy *SpotGridMartinStrategy) IsTrading() bool {
	return strategy.Trading
}

func (strategy *SpotGridMartinStrategy) IsInited() bool {
	return strategy.Inited
}
func (strategy *SpotGridMartinStrategy) Name() string {
	return strategy.StName
}

func (strategy *SpotGridMartinStrategy) Init(apiConfig *config.APIConfig) {
	strategy.pubWS = okex.NewOKExWSClient(okex.TEST_PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
	strategy.pubWS.ConnectWS()
	strategy.privWS = okex.NewOKExWSClient(okex.TEST_PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
	strategy.privWS.ConnectWS()

	symbols := []string{"ETH-USDT"}
	strategy.privWS.Login(apiConfig)
	strategy.pubWS.WatchTicker(symbols)
	strategy.privWS.WatchBalAndPos()
	strategy.pubWS.WatchDepth(symbols)
	strategy.Inited = true
	log.Infof("[%s] inited", strategy.StName)
}

func (strategy *SpotGridMartinStrategy) Symbol() string {
	return strategy.ssymbol
}

func (strategy *SpotGridMartinStrategy) OnTicker(tickerData bd.TickerData) {
	if !strategy.Inited {
		log.Infof("[%s] not init", strategy.Name())
		return
	}
	log.Infof("[SpotGridMartinStrategy]")
	//
}
