package app

import (
	"scindapsus/exchanges"
	"testing"
	"time"
)

func TestApp_Builder(t *testing.T) {
	ap := App{}
	builder := ap.Builder()
	builder.Exchange = exchanges.NewOKExchange()
	builder.Init()
	builder.Start()
	builder.Run()
}

func TestQiXianTaoLi(t *testing.T) {
	exchange := exchanges.NewOKExchange()
	//以ETH-USDT 为例
	//现货就是ETH-USDT,永续就是ETH-USDT-SWAP
	exchange.Init()
	exchange.WatchTicker("ETH-USDT")
	exchange.WatchTicker("ETH-USD-SWAP")

	time.Sleep(10 * time.Second)
}

func TestWatchFundingRate(t *testing.T) {
	exchange := exchanges.NewOKExchange()
	exchange.Init()
	exchange.WatchFundingRate("ETH-USDT-SWAP")
	time.Sleep(10 * time.Second)
}
