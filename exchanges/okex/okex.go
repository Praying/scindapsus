package okex

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"scindapsus/event"
	"scindapsus/websocket"
	"sync"
	"time"
)

type Exchange interface {
	SendLimitOrder()
}

func okexRespHandler(channel string, data json.RawMessage) error {
	switch channel {
	case "tickers":
		tickerData := parseTickerData(data)
		event.GetEventEngine().TickerChan <- (*tickerData)
		return nil
	default:
		return nil
	}
}

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
	okexWSClient.WsBuilder = websocket.NewWsBuilder().WsUrl(url).ReconnectInterval(time.Second).AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).ProtoHandleFunc(okexWSClient.handle)
	return okexWSClient
}

type wsResp struct {
	Event string `json:"event"`
	Arg   struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
	} `json:"arg"`
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (wsClient *OKExWSClient) handle(msg []byte) error {
	//TODO
	log.Debug("[ws][response]", string(msg))
	if string(msg) == "pong" {
		return nil
	}
	var wsResp wsResp
	err := json.Unmarshal(msg, &wsResp)
	if err != nil {
		log.Error(err)
		return err
	}

	if wsResp.Event != "" {
		switch wsResp.Event {
		case "subscribe":
			log.Info("subscribed:", wsResp.Arg.Channel)
			return nil
		case "unsubscribe":
			log.Info("unsubscribed:", wsResp.Arg.Channel)
			return nil
		case "error":
			log.Errorf(string(msg))
		default:
			log.Info(string(msg))
		}
		return fmt.Errorf("unknown websocket message: %v", wsResp)
	}

	if wsResp.Arg.Channel != "" {
		return wsClient.respHandler(wsResp.Arg.Channel, msg)
	}

	return nil
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
