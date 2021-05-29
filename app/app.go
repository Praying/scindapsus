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
	event.GetEventEngine().Init()
	strategy.GetStrategyEngine().Init()
}

func (this *AppBuilder) Run() {

	select {
	case <-this.DoneCh:

	}

	log.Infof("App ready to stop")

}
