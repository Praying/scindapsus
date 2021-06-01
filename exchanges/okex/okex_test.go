package okex

import (
	"encoding/json"
	"scindapsus/event"
	"testing"
	"time"
)

const PUBLIC_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/public"

const PRIVATE_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/private"

func TestSubscribeTickerData(t *testing.T) {
	event.GetEventEngine().Init()
	okexWSPublic := NewOKExWSClient(PUBLIC_WEBSOCKET_HOST_CHINA, okexRespHandler)
	subParam := json.RawMessage(`{ 
    "op": "subscribe",
    "args": [{
        "channel": "tickers",
        "instId": "BTC-USDT"
    }]}`)
	//okexWSPublic.Subscribe()
	okexWSPublic.ConnectWS()
	okexWSPublic.WSConn.SendMessage(subParam)
	time.Sleep(time.Second * 10)
}
