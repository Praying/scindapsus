package app

import (
	log "github.com/sirupsen/logrus"
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
	DoneCh         chan bool
}

func (this *App) Builder() AppBuilder {
	return AppBuilder{
		EventEngine: event.NewEventEngine(),
		DoneCh:      make(chan bool),
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

	select {
	case <-this.DoneCh:

	}

	log.Infof("App ready to stop")

}
