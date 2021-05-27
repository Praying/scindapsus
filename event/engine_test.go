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
			time.Sleep(time.Second * 1)
		}
	}(eventEngine)
	time.Sleep(time.Second * 5)
}
