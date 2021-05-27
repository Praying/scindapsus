package app

import (
	"scindapsus/event"
	"scindapsus/exchanges"
	"scindapsus/strategy"
)

type App struct {
}

type AppBuilder struct {
	EventEngine    *event.EventEngine
	Exchange       exchanges.Exchange
	StrategyEngine *strategy.StrategyEngine
}

func (this *App) Builder() AppBuilder {
	return AppBuilder{
		EventEngine: event.NewEventEngine(),
	}
}

func (this *AppBuilder) Init() {
	if this.Exchange != nil {
		this.Exchange.Init()
	}
	if this.EventEngine != nil {
		this.EventEngine.Init()
	}
	if this.StrategyEngine != nil {
		this.StrategyEngine.Init()
	}
}

func (this *AppBuilder) Run() {

}
