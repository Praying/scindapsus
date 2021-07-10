package strategy

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"scindapsus/config"
	"scindapsus/exchanges"
	"strings"
	"time"
)

type AbstractStrategy interface {
	Symbol() string
	OnTicker(tickerData bd.TickerData)
	OnTrade(tradeData bd.TradeData)
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
	Exchange exchanges.Exchange

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
	syms := strings.Split(strategy.ssymbol, "-")
	if len(syms) != 2 {
		log.Errorf("[%s] not correct ssymbol ", strategy.Name(), strategy.ssymbol)
		return
	}
	if bal, ok := balAndPosData.BalMap[syms[1]]; ok {
		strategy.balance = bal
		log.Infof("[%s] update balance: %v", strategy.Name(), strategy.balance)
	}
	//对于现货，持仓也属于余额，所以策略应该有个属性判断是期货还是现货
	if pos, ok := balAndPosData.BalMap[syms[0]]; ok {
		strategy.position = pos
		log.Infof("[%s] update position: %v", strategy.Name(), strategy.position)
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

//TODO 策略应该持有Exchange的引用
func (strategy *SpotGridMartinStrategy) Init(exchange exchanges.Exchange) {
	strategy.Exchange = exchange
	strategy.Exchange.WatchTicker("ETH-USDT", "", nil)
	strategy.Exchange.WatchBalance(nil)
	strategy.Exchange.WatchOrders("ETH-USDT", "", "", nil)
	strategy.Inited = true
	/*
		strategy.pubWS = okex.NewOKExWSClient(okex.TEST_PUBLIC_WEBSOCKET_HOST, okex.OkexRespHandler)
		strategy.pubWS.ConnectWS()
		strategy.privWS = okex.NewOKExWSClient(okex.TEST_PRIVATE_WEBSOCKET_HOST, okex.OkexRespHandler)
		strategy.privWS.ConnectWS()

		symbols := []string{"ETH-USDT"}
		strategy.privWS.Login(apiConfig)
		time.Sleep(1 * time.Second)
		strategy.pubWS.WatchTicker(symbols)
		strategy.privWS.WatchBalAndPos()
		strategy.privWS.WatchOrders(okex.MARGIN, symbols[0])
		strategy.pubWS.WatchDepth(symbols)
		strategy.Inited = true
	*/

	log.Infof("[%s] inited", strategy.StName)
}

func (strategy *SpotGridMartinStrategy) Symbol() string {
	return strategy.ssymbol
}

func (strategy *SpotGridMartinStrategy) OnTrade(tradeData bd.TradeData) {
	if tradeData.Direction == bd.Direction_LONG {

	}
	return
}

func (strategy *SpotGridMartinStrategy) OnTicker(tickerData bd.TickerData) {
	if !strategy.Inited {
		log.Infof("[%s] not init", strategy.Name())
		return
	}
	log.Infof("[SpotGridMartinStrategy] process ticker data, current price: %f", tickerData.Last)
	lastPrice := tickerData.Last
	if lastPrice > 2000 && strategy.position > 0 {
		//卖出
		//strategy.privWS.WatchCreateOrder(strategy.ssymbol, "limit", "sell", strategy.position, lastPrice)
		log.Infof("[%s] sell %s on price:%f, volume:%f", strategy.Name(), strategy.ssymbol, lastPrice, strategy.position)
	} else if strategy.position == 0 {
		//买入
		//strategy.privWS.WatchCreateOrder(strategy.ssymbol, "limit", "buy", 0.2, lastPrice)
		log.Infof("[%s] buy %s on price:%f, volume:%f", strategy.Name(), strategy.ssymbol, lastPrice, 0.2)
	}
	//
}
