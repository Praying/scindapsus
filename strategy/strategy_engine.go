package strategy

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

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

type TickerData struct {
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
	OrderType_RFQ    OrderType = "询价"
)

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
	Direction int32
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

	Volume    float64
	Frozen    float64
	Price     float64
	PNL       float64
	YD_Volume float64
}

type StrategyEngine struct {
	Strategies         []AbstractStrategy
	SymbolStrategyMap  map[string][]AbstractStrategy
	OrderIDStrategyMap map[string]AbstractStrategy
}

var once sync.Once
var instance *StrategyEngine

func GetStrategyEngine() *StrategyEngine {
	once.Do(func() {
		instance = NewStrategyEngine()
	})
	return instance
}

func NewStrategyEngine() *StrategyEngine {
	return &StrategyEngine{Strategies: make([]AbstractStrategy, 0)}
}

func (this *StrategyEngine) Init() {
	this.LoadStrategy()
	this.LoadStrategySetting()
	this.LoadStrategyData()
	log.Info("StrategyEngine init successfully")

}

func (this *StrategyEngine) LoadStrategy() {

}

func (this *StrategyEngine) LoadStrategySetting() {

}
func (this *StrategyEngine) LoadStrategyData() {

}
func (this *StrategyEngine) InitAllStrategies()  {}
func (this *StrategyEngine) StartAllStrategies() {}
func (this *StrategyEngine) ProcessTickerData(tickerData TickerData) {
	println("process ticker data")
}

func (this *StrategyEngine) ProcessBarData(barData BarData) {
	println("process bar data")
}

func (this *StrategyEngine) ProcessOrderData(orderData OrderData) {
	println("process order data")
}

func (this *StrategyEngine) ProcessTradeData(tradeData TradeData) {
	println("process trade data")
}

func (this *StrategyEngine) ProcessPositionData(positionData PositionData) {
	println("process postion data")
}
