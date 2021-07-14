package okex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"scindapsus/config"
	"scindapsus/websocket"
	"sync"
	"time"

	_ "scindapsus/config"
)

/*
产品类型
SPOT：币币
INST_MARGIN：币币杠杆
SWAP：永续合约
FUTURES：交割合约
OPTION：期权
ANY：全部
*/
const (
	INST_SPOT    string = "SPOT"        //币币
	INST_MARGIN  string = "INST_MARGIN" //币币杠杆
	INST_SWAP    string = "SWAP"        //永续合约
	INST_FUTURES string = "FUTURES"     //交割合约
	INST_OPTION  string = "OPTION"      //期权
	INST_ANY     string = "ANY"         //全部

	/*
		交易模式
		保证金模式 isolated：逐仓 cross： 全仓
		非保证金模式 cash：现金
	*/
	TRADE_MODEL_ISOLATED string = "isolated"
	TRADE_MODEL_CROSS    string = "cross"
	TRADE_MODEL_CASH     string = "cash"
	/*
		订单类型
		market：市价单
		limit：限价单
		post_only：只做maker单
		fok：全部成交或立即取消
		ioc：立即成交并取消剩余
		optimal_limit_ioc：市价委托立即成交并取消剩余（仅适用交割、永续）
	*/
	OKEX_OT_MARKET            string = "market"
	OKEX_OT_LIMIT             string = "limit"
	OKEX_OT_POST_ONLY         string = "post_only"
	OKEX_OT_FOK               string = "fok"
	OKEX_OT_IOC               string = "ioc"
	OKEX_OT_OPTIMAL_LIMIT_IOC string = "optimal_limit_ioc"
)

const (
	WS_SUBSCRIBE   string = "subscribe"
	WS_UNSUBSCRIBE string = "unsubscribe"
)

type APIConfig = config.APIConfig

const REST_HOST_CHINA string = "https://www.ouyi.cc"
const PUBLIC_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/public"
const PRIVATE_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/private"
const REST_HOST string = "https://www.okex.com"
const PUBLIC_WEBSOCKET_HOST string = "wss://ws.okex.com:8443/ws/v5/public"
const PRIVATE_WEBSOCKET_HOST string = "wss://ws.okex.com:8443/ws/v5/private"

const TEST_PUBLIC_WEBSOCKET_HOST string = "wss://wspri.coinall.ltd:8443/ws/v5/public?brokerId=9999"
const TEST_PRIVATE_WEBSOCKET_HOST string = "wss://wspri.coinall.ltd:8443/ws/v5/private?brokerId=9999"

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

func NewPublicWSClient(url string, respHandler func(channel string, data json.RawMessage) error) *PublicWSClient {
	publicWSClient := &PublicWSClient{
		once:        new(sync.Once),
		respHandler: respHandler,
	}
	publicWSClient.WsBuilder = websocket.NewWsBuilder().WsUrl(url).ReconnectInterval(time.Second).AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).ProtoHandleFunc(publicWSClient.handle)
	return publicWSClient
}

func NewPrivateWSClient(url string, respHandler func(channel string, data json.RawMessage) error) *PrivateWSClient {
	privateWSClient := &PublicWSClient{
		once:        new(sync.Once),
		respHandler: respHandler,
	}
	privateWSClient.WsBuilder = websocket.NewWsBuilder().WsUrl(url).ReconnectInterval(time.Second).AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).ProtoHandleFunc(privateWSClient.handle)
	return privateWSClient
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

func (wsClient *OKExWSClient) Subscribe(sub map[string]interface{}) error {
	wsClient.ConnectWS()
	return wsClient.WSConn.Subscribe(sub)
}

func generateOrderID(symbol string, orderType string) string {
	return fmt.Sprintf("%s-%s-%d", symbol, orderType, time.Now().Nanosecond())
}

func (wsClient *OKExWSClient) ConnectWS() {
	wsClient.once.Do(func() {
		wsClient.WSConn = wsClient.WsBuilder.Build()
	})
}

func genSign(secretKey string, timestamp string) string {
	msg := timestamp + "GET" + "/users/self/verify"
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type OKExRestClient struct {
}

type OKExExchange struct {
	//Rest和WebSocket
	publicWS   *PublicWSClient
	privateWS  *PrivateWSClient
	restClient *OKExRestClient
	config     *APIConfig
}

func NewOKExExchange(config *APIConfig) *OKExExchange {
	return &OKExExchange{config: config}
}
