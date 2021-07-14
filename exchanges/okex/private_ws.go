package okex

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"scindapsus/util"
	"strconv"
	"time"
)

type PrivateWSClient = OKExWSClient

func (wsClient *PrivateWSClient) WatchCreateOrder(symbol, orderType, side string, amount, price float64, clOrdID string) error {
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
		TdMode    string `json:"tdMode"`
		OrdType   string `json:"ordType"`
		/*
			当type为limit时，表示买入或卖出的数量
			当type为market时，现货交易买入时，表示买入的总金额，而
			当其他产品买入或卖出时，表示数量
		*/
		Sz string `json:"sz"`
		Px string `json:"px"`
	}{ClOrderID: clOrdID, Side: side, InstID: symbol, TdMode: TRADE_MODEL_CASH, OrdType: orderType, Sz: fmt.Sprintf("%f", amount), Px: fmt.Sprintf("%f", price)})
	data, err := json.Marshal(orderParam)
	if err != nil {
		log.Errorf("[ws][%s] json encode orderParam error , %s", wsClient.WSConn.WsUrl, err)
		return err
	}
	log.Info(string(data))
	wsClient.WSConn.SendMessage(data)
	return nil
}

//登录
func (wsClient *PrivateWSClient) Login(config *APIConfig) error {
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

//通过WebSocket撤单
//订单ID
//ordId和clOrdId必须传一个，若传两个，以 ordId 为主
func (wsClient *PrivateWSClient) WatchCancelOrder(orderId string, symbol string) {
	var cancelOrderParam CancelOrderParam
	//消息的唯一标识
	//用户提供，返回参数中会返回以便于找到相应的请求。
	//字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度必须要在1-32位之间。
	cancelOrderParam.ID = util.RandStringBytesMaskImprSrc(24)
	cancelOrderParam.Op = "cancel-order"
	cancelOrderParam.Args = append(cancelOrderParam.Args, struct {
		InstID  string `json:"instId"`
		OrdID   string `json:"ordId"`
		ClOrdID string `json:"clOrdId"`
	}{InstID: symbol, OrdID: orderId, ClOrdID: ""})
	wsClient.WSConn.Subscribe(cancelOrderParam)
}

func (wsClient *PrivateWSClient) WatchBalAndPos() error {
	return wsClient.doBalAndPos(WS_SUBSCRIBE)
}

func (wsClient *PrivateWSClient) UnWatchBalAndPos() error {
	return wsClient.doBalAndPos(WS_UNSUBSCRIBE)
}

func (wsClient *PrivateWSClient) doBalAndPos(op string) error {
	var param BalAndPosParam
	param.Op = op
	param.Args = append(param.Args, struct {
		Channel string `json:"channel"`
	}{Channel: BAL_AND_POS_CHANNEL})
	wsClient.WSConn.Subscribe(param)
	return nil
}

//订阅持仓
/*
	产品类型
INST_MARGIN：币币杠杆
SWAP：永续合约
FUTURES：交割合约
OPTION：期权
ANY：全部
*/

func (wsClient *PrivateWSClient) WatchPosition(instType string, instID string) error {
	return wsClient.doPosition(WS_SUBSCRIBE, instType, instID)
}
func (wsClient *PrivateWSClient) UnWatchPosition(instType string, instID string) error {
	return wsClient.doPosition(WS_UNSUBSCRIBE, instType, instID)
}

func (wsClient *PrivateWSClient) doPosition(op string, instType string, instID string) error {
	var positionParam PositionParam
	positionParam.Op = op
	positionParam.Args = append(positionParam.Args, struct {
		Channel  string `json:"channel"`  //必填
		InstType string `json:"instType"` //必填
		Uly      string `json:"uly"`
		InstID   string `json:"instId"`
	}{Channel: POSITIONS_CHANNEL, InstType: instType, Uly: "", InstID: instID})
	wsClient.WSConn.Subscribe(positionParam)
	return nil
}

//订单频道的订阅

func (wsClient *PrivateWSClient) WatchOrders(instType string, symbol string) error {
	return wsClient.doOrders(WS_SUBSCRIBE, instType, symbol)
}

func (wsClient *PrivateWSClient) UnWatchOrders(instType string, symbol string) error {
	return wsClient.doOrders(WS_UNSUBSCRIBE, instType, symbol)
}

func (wsClient *PrivateWSClient) doOrders(op string, instType string, symbol string) error {
	var orderChParam OrderChParam
	orderChParam.Op = op
	orderChParam.Args = append(orderChParam.Args, struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
		Uly      string `json:"uly"`
		InstID   string `json:"instId"`
	}{Channel: ORDERS_CHANNEL, InstType: instType, Uly: "", InstID: symbol})
	wsClient.WSConn.Subscribe(orderChParam)
	return nil
}
