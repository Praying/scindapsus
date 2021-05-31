package okex

import (
	"encoding/json"
	"net/http"
	"scindapsus/websocket"
	"sync"
	"time"
)

type Exchange interface {
	SendLimitOrder()
}

/*

type wsResp struct {
	Event     string `json:"event"`
	Channel   string `json:"channel"`
	Table     string `json:"table"`
	Data      json.RawMessage
	Success   bool        `json:"success"`
	ErrorCode interface{} `json:"errorCode"`
}

type OKExV3Ws struct {
	base *OKEx
	*WsBuilder
	once       *sync.Once
	WsConn     *WsConn
	respHandle func(channel string, data json.RawMessage) error
}

*/
type OKExWSClient struct {
	*websocket.WsBuilder
	config      *APIConfig
	once        *sync.Once
	WSConn      *websocket.WsConn
	respHandler func(channel string, data json.RawMessage) error
}

func NewOKExWSClient(url string, respHandler func(channel string, data json.RawMessage) error) *OKExWSClient {
	okexWSClient := &OKExWSClient{
		once:        new(sync.Once),
		respHandler: respHandler,
	}
	//okexWSClient.WsBuilder = websocket.NewWsBuilder().WsUrl(url).ReconnectInterval(time.Second).AutoReconnect().
	//	Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).ProtoHandleFunc(okexWSClient.respHandler)
}

/*
func NewOKExV3Ws(base *OKEx, handle func(channel string, data json.RawMessage) error) *OKExV3Ws {
	okV3Ws := &OKExV3Ws{
		once:       new(sync.Once),
		base:       base,
		respHandle: handle,
	}
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:8443/ws/v3").
		ReconnectInterval(time.Second).
		AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).
		DecompressFunc(FlateDecompress).ProtoHandleFunc(okV3Ws.handle)
	return okV3Ws
}

*/

type OKExRestClient struct {
}

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro

	Lever float64 //杠杆倍数 , for future
}
type OKExExchange struct {
	//Rest和WebSocket
	publicWS   *OKExWSClient
	privateWS  *OKExWSClient
	restClient *OKExRestClient
	config     *APIConfig
}

func NewOKExExchange(config *APIConfig) *OKExExchange {
	return &OKExExchange{config: config}
}
