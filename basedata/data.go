package basedata

import "time"

type Direction int32

const (
	Direction_SHORT Direction = iota
	Direction_LONG
	Direction_NET
)

type Interval string

const (
	Interval_MINUTE Interval = "1m"
	Interval_HOUR   Interval = "1h"
	Interval_DAILY  Interval = "1day"
	Interval_WEEKLY Interval = "1w"
	Interval_TICK   Interval = "tick"
)

const (
	Exchange_OKEX string = "OKEx"
)

type TickerData struct {
	Symbol        string    //货币符号
	TimeStamp     int64     //时间戳
	DateTime      time.Time //ISO8601
	High          float64   //	24小时最高价
	Low           float64   //24小时最低价
	Open          float64   //24小时开盘价
	Last          float64
	LastVolume    float64
	Close         float64
	PreviousClose float64
	Bid           float64 //买一价
	BidVolume     float64 //买一价的挂单数量
	Ask           float64 //卖一价
	AskVolume     float64 //卖一价的挂单数量
	BaseVolume    float64
	QuoteVolume   float64
	Percentage    float64
	Average       float64
	Change        float64
	VWap          float64
}

type BookRecord struct {
	Price  float64
	Volume float64
}

type BookData struct {
	Symbol    string
	TimeStamp int64
	Action    string
	AskList   []BookRecord
	BidList   []BookRecord
}

//deprecated
type TickerData2 struct {
	STSymbol string
	Symbol   string //货币符号
	Exchange string //交易所名称
	DateTime time.Time

	Volume    float64 //成交量
	OpenPrice float64 //开盘价
	HighPrice float64 //最高价
	LowPrice  float64 //最低价
	PreClose  float64

	//买一
	BidPrice1 float64
	BidPrice2 float64
	BidPrice3 float64
	BidPrice4 float64
	BidPrice5 float64

	//卖一
	AskPrice1 float64
	AskPrice2 float64
	AskPrice3 float64
	AskPrice4 float64
	AskPrice5 float64

	BidVolume1 float64
	BidVolume2 float64
	BidVolume3 float64
	BidVolume4 float64
	BidVolume5 float64

	AskVolume1 float64
	AskVolume2 float64
	AskVolume3 float64
	AskVolume4 float64
	AskVolume5 float64
}

type BarData struct {
	STSymbol string
	Symbol   string //货币符号
	Exchange string //交易所名称
	DateTime time.Time

	interval     Interval
	Volume       float64
	OpenInterest float64
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	ClosePrice   float64
}

type OrderType string

const (
	OrderType_LIMIT  OrderType = "限价"
	OrderType_MARKET OrderType = "市价"
	OrderType_STOP   OrderType = "STOP"
	OrderType_FAK    OrderType = "FAK"
	OrderType_FOK    OrderType = "FOK"
	OrderType_RFQ    OrderType = "询价"
)

var OrderType_Okex2Vt map[string]OrderType
var Direction_Okex2Vt map[string]Direction

//委托状态
var Status_Okex2Vt map[string]Status

func init() {
	OrderType_Okex2Vt = make(map[string]OrderType)
	OrderType_Okex2Vt["limit"] = OrderType_LIMIT
	OrderType_Okex2Vt["fok"] = OrderType_FOK

	Direction_Okex2Vt = make(map[string]Direction)
	Direction_Okex2Vt["buy"] = Direction_LONG
	Direction_Okex2Vt["sell"] = Direction_SHORT

	Status_Okex2Vt = make(map[string]Status)
	Status_Okex2Vt["live"] = Status_NOTTRADED
	Status_Okex2Vt["partially_filled"] = Status_PARTTRADED
	Status_Okex2Vt["filled"] = Status_ALLTRADED
	Status_Okex2Vt["canceled"] = Status_CANCELLED
}

type Offset string

const (
	Offset_NONE           Offset = ""
	Offset_OPEN           Offset = ""
	Offset_CLOSE          Offset = ""
	Offset_CLOSETODAY     Offset = ""
	Offset_CLOSEYESTERDAY Offset = ""
)

type Status string

const (
	Status_SUBMITTING Status = "提交中"
	Status_NOTTRADED  Status = "未成交"
	Status_PARTTRADED Status = "部分成交"
	Status_ALLTRADED  Status = "全部成交"
	Status_CANCELLED  Status = "已撤销"
	Status_REJECT     Status = "拒单"
)

type OrderData struct {
	Symbol    string //货币符号
	Exchange  string //交易所名称
	OrderID   string
	OrderType OrderType
	Direction Direction
	Offset    Offset
	Price     float64
	Volume    float64
	Traded    float64
	Status    Status
	DateTime  time.Time
	Reference string //这个字段是干啥的？
}

func (orderData *OrderData) IsActive() bool {
	return orderData.Status == Status_SUBMITTING || orderData.Status == Status_NOTTRADED || orderData.Status == Status_PARTTRADED
}

type TradeData struct {
	Symbol    string
	Exchange  string
	OrderID   string
	TradeID   string
	Direction Direction
	//Offset
	Price    float64
	Volume   float64
	Datetime time.Time
	//ST==Scindapasus Trader
	STSymbol  string
	STOrderID string
	STTradeID string
}

type PositionData struct {
	Symbol    string
	Exchange  string
	Direction int32

	//持仓数量
	Pos       float64
	Frozen    float64
	Price     float64
	PNL       float64
	YD_Volume float64
}

type BalAndPosData struct {
	//{"USDT":190.0, "ETH":312.3}
	BalMap      map[string]float64
	PositionMap map[string]float64
}

type FundingRateData struct {
	FundingRate     float64
	NextFundingRate float64
	FundingTime     time.Time
}

const (
	SIDE_BUY  string = "buy"
	SIDE_SELL string = "sell"
)
