package strategy

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"sync"
)

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
	return &StrategyEngine{Strategies: make([]AbstractStrategy, 0),
		SymbolStrategyMap: make(map[string][]AbstractStrategy, 0),
	}
}

func (this *StrategyEngine) Init() {
	this.LoadStrategy()
	this.LoadStrategySetting()
	this.LoadStrategyData()
	this.InitAllStrategies()
	log.Info("StrategyEngine init successfully")

}

func (this *StrategyEngine) LoadStrategy() {
	this.addStrategy(NewSpotGridMartinStrategy("ETH-USDT"))
}

func (this *StrategyEngine) LoadStrategySetting() {

}
func (this *StrategyEngine) LoadStrategyData() {

}
func (this *StrategyEngine) addStrategy(strategy AbstractStrategy) {
	this.Strategies = append(this.Strategies, strategy)
	if strategy == nil {
		log.Info("add empty strategy")
		return
	}
	this.SymbolStrategyMap[(strategy).Symbol()] = append(this.SymbolStrategyMap[strategy.Symbol()], strategy)
}
func (this *StrategyEngine) InitAllStrategies() {
	for _, strategy := range this.Strategies {
		strategy.Init()
	}
}
func (this *StrategyEngine) StartAllStrategies() {}
func (this *StrategyEngine) ProcessTickerData(tickerData bd.TickerData) {
	if len(this.SymbolStrategyMap[tickerData.Symbol]) > 0 {
		for _, strategy := range this.SymbolStrategyMap[tickerData.Symbol] {
			strategy.OnTicker(tickerData)
		}
	}
	log.Infof("process ticker data:%s", tickerData.Symbol)
}

func (this *StrategyEngine) ProcessBarData(barData bd.BarData) {
	log.Info("process bar data")
}

func (this *StrategyEngine) ProcessOrderData(orderData bd.OrderData) {
	log.Info("process order data")
}

func (this *StrategyEngine) ProcessTradeData(tradeData bd.TradeData) {
	log.Info("process trade data")
}

func (this *StrategyEngine) ProcessPositionData(positionData bd.PositionData) {
	log.Info("process postion data")
}

func (this *StrategyEngine) ProcessBookData(bookData bd.BookData) {
	log.Info("process book data: %v", bookData)
}
