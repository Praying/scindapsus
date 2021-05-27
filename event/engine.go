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
)

var Event2String map[EventType]string

func init() {
	Event2String = make(map[EventType]string)
	Event2String[Event_TICKER] = "event_ticker"
	Event2String[Event_ORDER] = "event_order"
	Event2String[Event_TRADE] = "event_trade"
	Event2String[Event_BAR] = "event_BAR"
}

type EventEngine struct {
	EventBus   evbus.Bus
	TickerChan chan strategy.TickerData
}

func NewEventEngine() *EventEngine {
	return &EventEngine{EventBus: evbus.New(),
		TickerChan: make(chan strategy.TickerData),
	}
}

func (eventEngine *EventEngine) Init() {
	eventEngine.EventBus.Subscribe(Event2String[Event_TICKER], func(tickerData strategy.TickerData) {
		strategy.GetInstance().ProcessTickerData(tickerData)
	})
	go func(eventEngine *EventEngine) {
		for {
			select {
			case tickerData := <-eventEngine.TickerChan:
				eventEngine.EventBus.Publish(Event2String[Event_TICKER], tickerData)
			}
		}
	}(eventEngine)
}
