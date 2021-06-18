package okex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"scindapsus/event"
	"scindapsus/websocket"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
const REST_HOST_CHINA: &str = "https://www.okex.win";
const PUBLIC_WEBSOCKET_HOST_CHINA: &str = "wss://wspri.coinall.ltd:8443/ws/v5/public";
const PRIVATE_WEBSOCKET_HOST_CHINA: &str = "wss://wspri.coinall.ltd:8443/ws/v5/private";

const REST_HOST: &str = "https://www.okex.com";
const PUBLIC_WEBSOCKET_HOST: &str = "wss://ws.okex.com:8443/ws/v5/public";
const PRIVATE_WEBSOCKET_HOST: &str = "wss://ws.okex.com:8443/ws/v5/private";

const TEST_PUBLIC_WEBSOCKET_HOST: &str = "wss://wspri.coinall.ltd:8443/ws/v5/public?brokerId=9999";
const TEST_PRIVATE_WEBSOCKET_HOST: &str =
    "wss://wspri.coinall.ltd:8443/ws/v5/private?brokerId=9999";
*/

const REST_HOST_CHINA string = "https://www.ouyi.cc"
const PUBLIC_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/public"
const PRIVATE_WEBSOCKET_HOST_CHINA string = "wss://wspri.coinall.ltd:8443/ws/v5/private"
const REST_HOST string = "https://www.okex.com"
const PUBLIC_WEBSOCKET_HOST string = "wss://ws.okex.com:8443/ws/v5/public"
const PRIVATE_WEBSOCKET_HOST string = "wss://ws.okex.com:8443/ws/v5/private"

const TEST_PUBLIC_WEBSOCKET_HOST string = "wss://wspri.coinall.ltd:8443/ws/v5/public?brokerId=9999"
const TEST_PRIVATE_WEBSOCKET_HOST string = "wss://wspri.coinall.ltd:8443/ws/v5/private?brokerId=9999"

type Exchange interface {
	CreateOrder()
}

func okexRespHandler(channel string, data json.RawMessage) error {
	switch channel {
	case "tickers":
		tickerData := parseTickerData(data)
		event.GetEventEngine().TickerChan <- (*tickerData)
		return nil
	case "books":
		fallthrough
	case "books5":
		fallthrough
	case "books-l2-tbt":
		fallthrough
	case "books50-l2-tbt":
		bookData := parseBookData(data)
		event.GetEventEngine().BookChan <- (*bookData)
		return nil
	case "position":
		positionData := parsePositionData(data)
		event.GetEventEngine().PositionChan <- (*positionData)
		return nil
	case "orders":
		orderData := parseOrderData(data)
		event.GetEventEngine().OrderChan <- (*orderData)
		return nil
	case "balance_and_position":
		log.Info(data)
		parseBalAndPosData(data)
		//event.GetEventEngine().BalAndPosChan
		return nil
	default:
		return nil
	}
	return nil
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

func (wsClient *OKExWSClient) Subscribe(sub map[string]interface{}) error {
	wsClient.ConnectWS()
	return wsClient.WSConn.Subscribe(sub)
}

func generateOrderID(symbol string, orderType string) string {
	return fmt.Sprintf("%s-%s-%d", symbol, orderType, time.Now().Nanosecond())
}
func (wsClient *OKExWSClient) CreateOrder(symbol, orderType, side string, amount, price float64) error {
	//Symbol-时间戳-OrderType
	//orderId := generateOrderID(symbol, orderType)
	orderId := fmt.Sprintf("%d", 1)
	orderParam := &OrderParam{
		ID:   orderId,
		Op:   "order",
		Args: nil,
	}

	orderParam.Args = append(orderParam.Args, struct {
		ClOrderID string `json:"clOrdId"`
		Side      string `json:"side"`
		InstID    string `json:"instId"`
		/*
			交易模式
			保证金模式 isolated：逐仓 cross： 全仓
			非保证金模式 cash：现金
		*/
		TdMode  string `json:"tdMode"`
		OrdType string `json:"ordType"`
		/*
			当type为limit时，表示买入或卖出的数量
			当type为market时，现货交易买入时，表示买入的总金额，而
			当其他产品买入或卖出时，表示数量
		*/
		Sz string `json:"sz"`
		Px string `json:"px"`
	}{ClOrderID: "xxsyydflsdfdsuf", Side: side, InstID: symbol, TdMode: "cash", OrdType: orderType, Sz: fmt.Sprintf("%f", amount), Px: fmt.Sprintf("%f", price)})
	data, err := json.Marshal(orderParam)
	if err != nil {
		log.Errorf("[ws][%s] json encode orderParam error , %s", wsClient.WSConn.WsUrl, err)
		return err
	}
	log.Info(string(data))
	wsClient.WSConn.SendMessage(data)
	return nil
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

//登录
func (wsClient *OKExWSClient) Login(config *APIConfig) error {
	var loginParam LoginParam
	loginParam.Op = "login"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	loginParam.Args = append(loginParam.Args, struct {
		APIKey     string `json:"apiKey"`
		Passphrase string `json:"passphrase"`
		Timestamp  string `json:"timestamp"`
		Sign       string `json:"sign"`
	}{APIKey: config.ApiKey, Passphrase: config.ApiPassphrase, Timestamp: timestamp, Sign: genSign(config.ApiSecretKey, timestamp)})
	data, err := json.Marshal(loginParam)
	if err != nil {
		log.Errorf("[ws][%s] json encode error , %s", wsClient.WSConn.WsUrl, err)
		return err
	}
	log.Debug(string(data))
	wsClient.WSConn.SendMessage(data)
	return nil
}

//订阅行情Ticker数据
func (wsClient *OKExWSClient) SubscribeTicker(currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = "subscribe"
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{"tickers", currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

func (wsClient *OKExWSClient) UnSubscribeTicker(currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = "unsubscribe"
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{"tickers", currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

//订阅持仓
func (wsClient *OKExWSClient) SubscribePosition(instType string, instID string) error {
	//TODO
	var positionParam PositionParam
	positionParam.Op = "subscribe"
	positionParam.Args = append(positionParam.Args, struct {
		Channel  string `json:"channel"`  //必填
		InstType string `json:"instType"` //必填
		Uly      string `json:"uly"`
		InstID   string `json:"instId"`
	}{Channel: "positions", InstType: instType, Uly: "", InstID: instID})
	wsClient.WSConn.Subscribe(positionParam)
	return nil
}

//订阅深度数据
func (wsClient *OKExWSClient) SubscribeDepth(currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = "subscribe"
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{"books", currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

func (wsClient *OKExWSClient) SubscribeBalAndPos() error {
	param := `{
    "op": "subscribe",
    "args": [{
        "channel": "balance_and_position"
    }]
	}`
	wsClient.WSConn.SendMessage([]byte(param))
	return nil
}

func (wsClient *OKExWSClient) UnSubscribeDepth(currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = "unsubscribe"
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{"books", currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

func (wsClient *OKExWSClient) handle(msg []byte) error {
	//TODO
	log.Info("[ws][response]", string(msg))
	if string(msg) == "pong" {
		return nil
	} else if strings.Contains(string(msg), "clOrdId") {
		//处理下单的状态
		var orderResp OrderResp
		err := json.Unmarshal(msg, &orderResp)
		if err != nil {
			log.Error(err)
			return err
		}
		if orderResp.Code == "1" {
			if len(orderResp.Data) > 0 {
				log.Errorf("clOrdId:%s, ordId:%s,sCode:%s, sMsg:%s", orderResp.Data[0].ClOrdID, orderResp.Data[0].OrdID, orderResp.Data[0].SCode, orderResp.Data[0].SMsg)
				return fmt.Errorf("clOrdId:%s, ordId:%s,sCode:%s, sMsg:%s", orderResp.Data[0].ClOrdID, orderResp.Data[0].OrdID, orderResp.Data[0].SCode, orderResp.Data[0].SMsg)
			}
		} else if orderResp.Code == "0" {
			log.Infof("order successful: ordId: %s", orderResp.Data[0].OrdID)
			return nil
		}

	} else if strings.Contains(string(msg), "event") {
		//处理订阅事件的状态
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
			case "login":
				log.Info("login:", string(msg))
			case "error":
				log.Errorf(string(msg))
			default:
				log.Info(string(msg))
			}
			return fmt.Errorf("unknown websocket message: %v", wsResp)
		}
		if wsResp.Code != "0" {
			log.Errorf("error")
		}

		if wsResp.Arg.Channel != "" {
			return wsClient.respHandler(wsResp.Arg.Channel, msg)
		}

	}

	return nil
}

//func (wsClient *OKExWSClient) SubscribePosition(s string, s2 string, s3 string, i int, i2 int) interface{} {
//
//}

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
