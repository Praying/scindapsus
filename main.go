package scindapsus

import (
	"scindapsus/app"
	"scindapsus/exchanges"
	"scindapsus/strategy"
)

func main() {
	ap := app.App{}
	builder := ap.Builder()
	builder.StrategyEngine = strategy.NewStrategyEngine()
	builder.Exchange = exchanges.NewOKExchange()
	builder.Init()
	builder.Run()
}
