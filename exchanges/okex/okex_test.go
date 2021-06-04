package okex

import (
	"github.com/stretchr/testify/assert"
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

func TestGenSign(t *testing.T) {
	sign := genSign("0C4D537A249467C8102EAE10E6118ED6", "1538054050")
	assert.Equal(t, "fzrMT2u+wBZ1oq++co4MuDTvjmKNWcL7PClJg0y2Qj4=", sign)
}

func TestPrivateWS(t *testing.T) {
	apiConfig := &APIConfig{
		HttpClient:    nil,
		Endpoint:      "",
		ApiKey:        "6b7c638e-0486-40d0-9228-db4d53f585f6",
		ApiSecretKey:  "B22D93123C8096E203D8A47252E277C8",
		ApiPassphrase: "cyjqr1314",
		ClientId:      "",
		Lever:         0,
	}
	privateWS := NewOKExWSClient(PRIVATE_WEBSOCKET_HOST_CHINA, okexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)
	//订阅持仓信息
	privateWS.SubscribePosition("MARGIN", "SHIB-USDT")
	time.Sleep(5 * time.Second)
}
