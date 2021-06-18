package okex

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"scindapsus/event"
	"testing"
	"time"
)

var apiConfig *APIConfig

func init() {
	apiConfig = &APIConfig{
		HttpClient:    nil,
		Endpoint:      "",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
		ClientId:      "",
		Lever:         0,
	}
}

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

	privateWS := NewOKExWSClient(PRIVATE_WEBSOCKET_HOST_CHINA, okexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)
	//订阅持仓信息
	privateWS.SubscribePosition("MARGIN", "SHIB-USDT")
	time.Sleep(5 * time.Second)
}

func TestSendLimitOrder(t *testing.T) {

	privateWS := NewOKExWSClient(TEST_PRIVATE_WEBSOCKET_HOST, okexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)
	//订阅订单信息

	//下单
	if err := privateWS.CreateOrder("ETH-USDT", "limit", "buy", 0.01, 2400); err != nil {
		log.Errorln(err.Error())
	}
	time.Sleep(5 * time.Second)
}

//订阅余额和持仓
func TestSubscribeBalAndPos(t *testing.T) {
	privateWS := NewOKExWSClient(TEST_PRIVATE_WEBSOCKET_HOST, okexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)

	if err := privateWS.SubscribeBalAndPos(); err != nil {
		log.Errorln(err.Error())
	}
	time.Sleep(5 * time.Second)
}
