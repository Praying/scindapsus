package event

import (
	evbus "github.com/asaskevich/EventBus"
	"scindapsus/strategy"
	_ "scindapsus/strategy"
)

type EventType int32

const (
	Event_TICKER EventType = iota
	Event_ORDER
	Event_TRADE
	Event_BAR
	Event_POSITION
)

const DEFUALT_CHANNEL_SIZE = 100

var Event2String map[EventType]string

func init() {
	Event2String = make(map[EventType]string)
	Event2String[Event_TICKER] = "event_ticker"
	Event2String[Event_ORDER] = "event_order"
	Event2String[Event_TRADE] = "event_trade"
	Event2String[Event_BAR] = "event_bar"
	Event2String[Event_POSITION] = "event_position"
}

type EventEngine struct {
	EventBus     evbus.Bus
	TickerChan   chan strategy.TickerData
	OrderChan    chan strategy.OrderData
	TradeChan    chan strategy.TradeData
	PositionChan chan strategy.PositionData
}

func NewEventEngine() *EventEngine {
	return &EventEngine{EventBus: evbus.New(),
		TickerChan:   make(chan strategy.TickerData, DEFUALT_CHANNEL_SIZE),
		OrderChan:    make(chan strategy.OrderData, DEFUALT_CHANNEL_SIZE),
		TradeChan:    make(chan strategy.TradeData, DEFUALT_CHANNEL_SIZE),
		PositionChan: make(chan strategy.PositionData, DEFUALT_CHANNEL_SIZE),
	}
}

func (eventEngine *EventEngine) Init() {
	eventEngine.EventBus.Subscribe(Event2String[Event_TICKER], func(tickerData strategy.TickerData) {
		strategy.GetInstance().ProcessTickerData(tickerData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_ORDER], func(orderData strategy.OrderData) {
		strategy.GetInstance().ProcessOrderData(orderData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_TRADE], func(tradeData strategy.TradeData) {
		strategy.GetInstance().ProcessTradeData(tradeData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_POSITION], func(positionData strategy.PositionData) {
		strategy.GetInstance().ProcessPositionData(positionData)
	})
	go func(eventEngine *EventEngine) {
		for {
			select {
			case tickerData := <-eventEngine.TickerChan:
				eventEngine.EventBus.Publish(Event2String[Event_TICKER], tickerData)
			case orderData := <-eventEngine.OrderChan:
				eventEngine.EventBus.Publish(Event2String[Event_ORDER], orderData)
			case tradeData := <-eventEngine.TradeChan:
				eventEngine.EventBus.Publish(Event2String[Event_TRADE], tradeData)
			case positionData := <-eventEngine.PositionChan:
				eventEngine.EventBus.Publish(Event2String[Event_POSITION], positionData)
			}
		}
	}(eventEngine)
}
