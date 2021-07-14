package okex

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	bd "scindapsus/basedata"
	"scindapsus/event"
	"scindapsus/util"
	"testing"
	"time"
)

var apiConfig *APIConfig

func init() {
	apiConfig = &APIConfig{
		HttpClient:    nil,
		Endpoint:      "",
		ApiKey:        "befbebe7-55af-4f86-a484-ed6eb8553879",
		ApiSecretKey:  "3F65A4EB0A0EEFB1E31DA35982E1B8E6",
		ApiPassphrase: "cyjqr1314",
		ClientId:      "",
		Lever:         0,
	}
}

func TestSubscribeTicker(t *testing.T) {
	event.GetEventEngine().Init()
	okexWSPublic := NewOKExWSClient(PUBLIC_WEBSOCKET_HOST_CHINA, OkexRespHandler)
	okexWSPublic.ConnectWS()
	okexWSPublic.WatchTicker([]string{"BTC-USDT", "ETH-USDT"})
	time.Sleep(time.Second * 5)
	okexWSPublic.UnWatchTicker([]string{"BTC-USDT"})
	time.Sleep(time.Second * 5)
}

func TestSubscribeDepth(t *testing.T) {
	event.GetEventEngine().Init()
	okexWSPublic := NewPublicWSClient(PUBLIC_WEBSOCKET_HOST_CHINA, OkexRespHandler)
	okexWSPublic.ConnectWS()
	okexWSPublic.WatchDepth([]string{"BTC-USDT", "ETH-USDT"})
	time.Sleep(time.Second * 5)
	okexWSPublic.UnWatchDepth([]string{"BTC-USDT"})
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

	privateWS := NewPrivateWSClient(PRIVATE_WEBSOCKET_HOST_CHINA, OkexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)
	//订阅持仓信息
	privateWS.WatchPosition(INST_SPOT, "SHIB-USDT")
	time.Sleep(5 * time.Second)
}

func TestSendLimitOrder(t *testing.T) {

	privateWS := NewPrivateWSClient(TEST_PRIVATE_WEBSOCKET_HOST, OkexRespHandler)
	privateWS.ConnectWS()
	//登录
	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)
	//订阅订单信息

	//下单
	if err := privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, bd.SIDE_BUY, 0.01, 2400, ""); err != nil {
		log.Errorln(err.Error())
	}
	time.Sleep(5 * time.Second)
}

//订阅余额和持仓
func TestSubscribeBalAndPos(t *testing.T) {
	privateWS := NewPrivateWSClient(TEST_PRIVATE_WEBSOCKET_HOST, OkexRespHandler)
	privateWS.ConnectWS()
	//登录

	privateWS.Login(apiConfig)
	time.Sleep(5 * time.Second)

	if err := privateWS.WatchBalAndPos(); err != nil {
		log.Errorln(err.Error())
	}
	privateWS.WatchPosition(INST_SPOT, "ETH-USDT")
	time.Sleep(5 * time.Second)
	//下单
	if err := privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, bd.SIDE_BUY, 0.01, 2400, ""); err != nil {
		log.Errorln(err.Error())
	}
	time.Sleep(2 * time.Second)
	if err := privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, bd.SIDE_BUY, 0.01, 2400, ""); err != nil {
		log.Errorln(err.Error())
	}
	time.Sleep(2 * time.Second)
}

//订阅订单信息
func TestOKExWSClient_WatchOrders(t *testing.T) {
	privateWS := NewPrivateWSClient(TEST_PRIVATE_WEBSOCKET_HOST, OkexRespHandler)
	privateWS.ConnectWS()
	//登录

	privateWS.Login(apiConfig)
	time.Sleep(2 * time.Second)
	privateWS.WatchOrders(INST_SPOT, "ETH-USDT")
	time.Sleep(2 * time.Second)
	ts := time.Now().Unix()
	clOrdID := util.GenerateClOrdId(ts, 1)
	privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, "buy", 0.2, 2050, clOrdID)
	time.Sleep(2 * time.Second)
	clOrdID = util.GenerateClOrdId(ts, 2)
	privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, "buy", 0.2, 2050, clOrdID)
	time.Sleep(2 * time.Second)
	clOrdID = util.GenerateClOrdId(ts, 3)
	privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, "buy", 0.2, 2050, clOrdID)
	time.Sleep(2 * time.Second)
	clOrdID = util.GenerateClOrdId(ts, 4)
	privateWS.WatchCreateOrder("ETH-USDT", OKEX_OT_LIMIT, "buy", 0.2, 2050, clOrdID)
	time.Sleep(2 * time.Second)

}
