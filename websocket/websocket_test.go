package websocket

import (
	"encoding/json"
	"testing"
	"time"
)

func ProtoHandle(data []byte) error {
	println(string(data))
	return nil
}

func TestNewConn(t *testing.T) {
	var heartbeatFunc = func() []byte {
		ts := time.Now().Unix()*1000 + 42029
		args := make([]interface{}, 0)
		args = append(args, ts)
		ping := "ping"

		return []byte(ping)
	}

	ws := NewWsBuilder().Dump().WsUrl("wss://wspri.coinall.ltd:8443/ws/v5/public").
		AutoReconnect().
		Heartbeat(heartbeatFunc, 5*time.Second).ProtoHandleFunc(ProtoHandle).Build()

	param := json.RawMessage(`	{
	    "op": "subscribe",
	    "args": [{
	        "channel": "tickers",
	        "instId": "DOGE-USDT"
	    }]
	}`)
	t.Log(ws.Subscribe(&param))
	/*

		{
		    "op": "subscribe",
		    "args": [{
		        "channel": "tickers",
		        "instId": "LTC-USD-200327"
		    }]
		}
	*/
	time.Sleep(time.Second * 20)
	ws.c.Close()
	time.Sleep(time.Second * 120)
}
