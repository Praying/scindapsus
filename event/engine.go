package event

import (
	evbus "github.com/asaskevich/EventBus"
	log "github.com/sirupsen/logrus"
	"scindapsus/strategy"
	_ "scindapsus/strategy"
	"sync"
)

type EventType int32

const (
	//行情
	Event_TICKER EventType = iota
	Event_ORDER
	Event_TRADE
	Event_BAR
	Event_POSITION
	//深度
	Event_BOOK
)

var once sync.Once
var instance *EventEngine

func GetEventEngine() *EventEngine {
	once.Do(func() {
		instance = NewEventEngine()
	})
	return instance
}

const DEFUALT_CHANNEL_SIZE = 100

var Event2String map[EventType]string

func init() {
	Event2String = make(map[EventType]string)
	Event2String[Event_TICKER] = "event_ticker"
	Event2String[Event_ORDER] = "event_order"
	Event2String[Event_TRADE] = "event_trade"
	Event2String[Event_BAR] = "event_bar"
	Event2String[Event_POSITION] = "event_position"
	Event2String[Event_BOOK] = "event_book"
}

type EventEngine struct {
	EventBus     evbus.Bus
	TickerChan   chan strategy.TickerData
	OrderChan    chan strategy.OrderData
	TradeChan    chan strategy.TradeData
	PositionChan chan strategy.PositionData
	BookChan     chan strategy.BookData
}

func NewEventEngine() *EventEngine {
	return &EventEngine{EventBus: evbus.New(),
		TickerChan:   make(chan strategy.TickerData, DEFUALT_CHANNEL_SIZE),
		OrderChan:    make(chan strategy.OrderData, DEFUALT_CHANNEL_SIZE),
		TradeChan:    make(chan strategy.TradeData, DEFUALT_CHANNEL_SIZE),
		PositionChan: make(chan strategy.PositionData, DEFUALT_CHANNEL_SIZE),
		BookChan:     make(chan strategy.BookData, DEFUALT_CHANNEL_SIZE),
	}
}

func (eventEngine *EventEngine) Init() {
	eventEngine.EventBus.Subscribe(Event2String[Event_TICKER], func(tickerData strategy.TickerData) {
		strategy.GetStrategyEngine().ProcessTickerData(tickerData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_ORDER], func(orderData strategy.OrderData) {
		strategy.GetStrategyEngine().ProcessOrderData(orderData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_TRADE], func(tradeData strategy.TradeData) {
		strategy.GetStrategyEngine().ProcessTradeData(tradeData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_POSITION], func(positionData strategy.PositionData) {
		strategy.GetStrategyEngine().ProcessPositionData(positionData)
	})
	eventEngine.EventBus.Subscribe(Event2String[Event_BOOK], func(bookData strategy.BookData) {
		strategy.GetStrategyEngine().ProcessBookData(bookData)
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
			case bookData := <-eventEngine.BookChan:
				eventEngine.EventBus.Publish(Event2String[Event_BOOK], bookData)
			}
		}
	}(eventEngine)
	log.Info("EventEngine init successfully")
}
