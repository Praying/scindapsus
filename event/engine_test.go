package event

import (
	bd "scindapsus/basedata"
	"testing"
	"time"
)

func TestProcessEvent(t *testing.T) {
	eventEngine := NewEventEngine()
	eventEngine.Init()
	go func(engine *EventEngine) {
		for i := 0; i < 3; i++ {
			tickerData := bd.TickerData{}
			engine.TickerChan <- tickerData
			orderData := bd.OrderData{}
			engine.OrdersChan <- orderData
			tradeData := bd.TradeData{}
			engine.TradeChan <- tradeData
			positionData := bd.PositionData{}
			engine.PositionChan <- positionData
			time.Sleep(time.Second * 1)
		}
	}(eventEngine)
	time.Sleep(time.Second * 5)
}
