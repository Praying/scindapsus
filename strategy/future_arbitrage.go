package strategy

import (
	log "github.com/sirupsen/logrus"
	"math"
	bd "scindapsus/basedata"
	"scindapsus/exchanges"
	"scindapsus/exchanges/okex"
)

type FutureArbitrage struct {
	//该策略可能涉及到多个交易对，如ETH-USDT, ETH-USD-SWAP
	SD_Symbols []string
	//该策略的名字
	SD_Name string

	//记录交易对的价格
	SD_SymbolPrice map[string]float64

	SD_IsInited bool

	SD_Exchange exchanges.Exchange

	//等待订单成交
	SD_Wait bool
}

func NewFutureArbitrage(SD_Symbols []string) *FutureArbitrage {
	return &FutureArbitrage{SD_Symbols: SD_Symbols}
}

func (f *FutureArbitrage) Symbol() string {
	return ""
}

func (f *FutureArbitrage) Symbols() []string {
	return f.SD_Symbols
}

func (f *FutureArbitrage) OnTicker(tickerData bd.TickerData) {
	if f.SD_Wait {
		return
	}
	//默认顺序是现货-永续，不可打乱
	f.SD_SymbolPrice[tickerData.Symbol] = tickerData.Last
	//计算价差，超过5%，现货买入，期货开空
	if _, ok := f.SD_SymbolPrice[f.SD_Symbols[0]]; !ok {
		return
	}
	if _, ok := f.SD_SymbolPrice[f.SD_Symbols[1]]; !ok {
		return
	}
	spotPrice := f.SD_SymbolPrice[f.SD_Symbols[0]]
	swapPrice := f.SD_SymbolPrice[f.SD_Symbols[1]]
	if spotPrice < 0.000001 || swapPrice < 0.000001 {
		return
	}
	diff := math.Abs(spotPrice-swapPrice) / spotPrice
	if diff > 0.05 {
		log.Infof("spot:%f, swap:%f, diff:%f%%", spotPrice, swapPrice, diff*100)
		//现货做多
		f.SD_Exchange.WatchCreateOrder(f.SD_Symbols[0], okex.OKEX_OT_IOC, bd.SIDE_BUY, 1, spotPrice, okex.TRADE_MODEL_CASH)
		//永续做空
		f.SD_Exchange.WatchCreateOrder(f.SD_Symbols[1], okex.OKEX_OT_LIMIT, bd.SIDE_SELL, 1, swapPrice, okex.TRADE_MODEL_CROSS)

		f.SD_Wait = true
	}
}

func (f FutureArbitrage) OnTrade(tradeData bd.TradeData) {
	log.Infof("%s, traded volume:%f, trade price:%f", tradeData.Symbol, tradeData.Volume, tradeData.Price)
}

func (f FutureArbitrage) OnPosition(positionData bd.PositionData) {
	log.Infof("%s, pos:%f ", positionData.Symbol, positionData.Pos)
}

func (f FutureArbitrage) OnBalAndPos(balAndPosData bd.BalAndPosData) {
	panic("implement me")
}

func (f *FutureArbitrage) Init(exchange exchanges.Exchange) {
	f.SD_SymbolPrice = make(map[string]float64)
	for _, symbol := range f.SD_Symbols {
		f.SD_SymbolPrice[symbol] = 0.0
	}

	f.SD_Exchange = exchange
	f.SD_Exchange.WatchInstruments([]string{okex.INST_SPOT, okex.INST_SWAP})
	for _, symbol := range f.SD_Symbols {
		f.SD_Exchange.WatchTicker(symbol)
		f.SD_Exchange.WatchPosition(symbol)
		f.SD_Exchange.WatchOrders(symbol, "", "", nil)
		//f.SD_Exchange.WatchTrades(symbol,"","",nil)
	}
	f.SD_IsInited = true

}

func (f *FutureArbitrage) Name() string {
	return f.SD_Name
}

func (f *FutureArbitrage) IsInited() bool {
	return f.SD_IsInited
}

func (f FutureArbitrage) OnStart() {
	panic("implement me")
}

func (f FutureArbitrage) IsTrading() bool {
	panic("implement me")
}

func (f FutureArbitrage) OnStop() {
	panic("implement me")
}
