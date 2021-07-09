package app

import (
	"scindapsus/exchanges"
	"testing"
)

func TestApp_Builder(t *testing.T) {
	ap := App{}
	builder := ap.Builder()
	builder.Exchange = exchanges.NewOKExchange()
	builder.Init()
	//symbol:="ETH-USDT"
	//builder.Exchange.WatchOrders(symbol, "", "", nil)
	//builder.Exchange.WatchPosition(symbol)
	//builder.Exchange.WatchBalance(nil)
	//builder.Exchange.WatchTicker(symbol,"",nil)
	builder.Start()

	builder.Run()
}

func TestOther(t *testing.T) {

}
