package event

import (
	"scindapsus/strategy"
	"testing"
	"time"
)

func TestProcessEvent(t *testing.T) {
	eventEngine := NewEventEngine()
	eventEngine.Init()
	go func(engine *EventEngine) {
		for i := 0; i < 3; i++ {
			tickerData := strategy.TickerData{}
			engine.TickerChan <- tickerData
			orderData := strategy.OrderData{}
			engine.OrderChan <- orderData
			tradeData := strategy.TradeData{}
			engine.TradeChan <- tradeData
			positionData := strategy.PositionData{}
			engine.PositionChan <- positionData
			time.Sleep(time.Second * 1)
		}
	}(eventEngine)
	time.Sleep(time.Second * 5)
}
