package okex

import (
	"scindapsus/event"
	"testing"
	"time"
)

const PUBLIC_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/public"

const PRIVATE_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/private"

func TestSubscribeTicker(t *testing.T) {
	event.GetEventEngine().Init()
	okexWSPublic := NewOKExWSClient(PUBLIC_WEBSOCKET_HOST_CHINA, okexRespHandler)
	okexWSPublic.ConnectWS()
	okexWSPublic.SubscribeTicker([]string{"BTC-USDT", "ETH-USDT"})
	time.Sleep(time.Second * 5)
	okexWSPublic.UnSubscribeTicker([]string{"BTC-USDT"})
	time.Sleep(time.Second * 5)
}

func TestSubscribeDepth(t *testing.T) {
	event.GetEventEngine().Init()
	okexWSPublic := NewOKExWSClient(PUBLIC_WEBSOCKET_HOST_CHINA, okexRespHandler)
	okexWSPublic.ConnectWS()
	okexWSPublic.SubscribeDepth([]string{"BTC-USDT", "ETH-USDT"})
	time.Sleep(time.Second * 5)
	okexWSPublic.UnSubscribeDepth([]string{"BTC-USDT"})
	time.Sleep(time.Second * 5)
}

func TestReflect(t *testing.T) {
	var s SubParam
	s.Args = append(s.Args, struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	}{"tickers", "BTC-USDT"})
	s.Args = append(s.Args, struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	}{"tickers", "ETH-USDT"})

}
