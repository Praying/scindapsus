package event

import (
	evbus "github.com/asaskevich/EventBus"
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"

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
	Event_BALANDPOS
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
	Event2String[Event_BALANDPOS] = "event_bal_and_pos"
}

type EventEngine struct {
	EventBus      evbus.Bus
	TickerChan    chan bd.TickerData
	OrderChan     chan bd.OrderData
	TradeChan     chan bd.TradeData
	PositionChan  chan bd.PositionData
	BookChan      chan bd.BookData
	BalAndPosChan chan bd.BalAndPosData
}

func NewEventEngine() *EventEngine {
	return &EventEngine{EventBus: evbus.New(),
		TickerChan:    make(chan bd.TickerData, DEFUALT_CHANNEL_SIZE),
		OrderChan:     make(chan bd.OrderData, DEFUALT_CHANNEL_SIZE),
		TradeChan:     make(chan bd.TradeData, DEFUALT_CHANNEL_SIZE),
		PositionChan:  make(chan bd.PositionData, DEFUALT_CHANNEL_SIZE),
		BookChan:      make(chan bd.BookData, DEFUALT_CHANNEL_SIZE),
		BalAndPosChan: make(chan bd.BalAndPosData, DEFUALT_CHANNEL_SIZE),
	}
}

func (eventEngine *EventEngine) Init() {

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
			case balAndPosData := <-eventEngine.BalAndPosChan:
				eventEngine.EventBus.Publish(Event2String[Event_BALANDPOS], balAndPosData)
			}
		}
	}(eventEngine)
	log.Info("EventEngine init successfully")
}
