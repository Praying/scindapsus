package app

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"scindapsus/event"
	"scindapsus/exchanges"
	"scindapsus/strategy"
)

type App struct {
}

type AppBuilder struct {
	Exchange exchanges.Exchange
	DoneCh   chan bool
}

func (this *App) Builder() AppBuilder {
	return AppBuilder{
		DoneCh: make(chan bool),
	}
}

func (this *AppBuilder) Init() {
	if this.Exchange != nil {
		this.Exchange.Init()
	}
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_TICKER], func(tickerData bd.TickerData) {
		strategy.GetStrategyEngine().ProcessTickerData(tickerData)
	})
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_ORDERS], func(orderData bd.OrderData) {
		strategy.GetStrategyEngine().ProcessOrdersData(orderData)
	})
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_TRADE], func(tradeData bd.TradeData) {
		strategy.GetStrategyEngine().ProcessTradeData(tradeData)
	})
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_POSITION], func(positionData bd.PositionData) {
		strategy.GetStrategyEngine().ProcessPositionData(positionData)
	})
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_BOOK], func(bookData bd.BookData) {
		strategy.GetStrategyEngine().ProcessBookData(bookData)
	})
	event.GetEventEngine().EventBus.Subscribe(event.Event2String[event.Event_BALANDPOS], func(balAndPosData bd.BalAndPosData) {
		strategy.GetStrategyEngine().ProcessBalAndPosData(balAndPosData)
	})
	event.GetEventEngine().Init()
	strategy.GetStrategyEngine().Init(this.Exchange)
}

func (this *AppBuilder) Run() {

	select {
	case <-this.DoneCh:

		break
	}

	log.Infof("App ready to stop")

}

func (this *AppBuilder) Start() {
	strategy.GetStrategyEngine().StartAllStrategies()
}
