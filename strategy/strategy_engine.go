package strategy

import (
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"
	"scindapsus/exchanges"
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

func (this *StrategyEngine) Init(exchange exchanges.Exchange) {
	this.LoadStrategy()
	this.LoadStrategySetting()
	this.LoadStrategyData()
	this.InitAllStrategies(exchange)
	log.Info("StrategyEngine init successfully")

}

func (this *StrategyEngine) LoadStrategy() {
	//this.addStrategy(NewSpotGridMartinStrategy("ETH-USDT"))
	this.addStrategy(NewFutureArbitrage([]string{"ETH-USDT", "ETH-USD-SWAP"}))
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
	for _, symbol := range strategy.Symbols() {
		this.SymbolStrategyMap[symbol] = append(this.SymbolStrategyMap[strategy.Symbol()], strategy)
	}

}
func (this *StrategyEngine) InitAllStrategies(exchange exchanges.Exchange) {
	for _, strategy := range this.Strategies {
		go func(strategy AbstractStrategy) {
			if strategy.IsInited() {
				log.Info("[%s] has been inited", strategy.Name())
				return
			}
			strategy.Init(exchange)
			//数据可能需要恢复
			//行情订阅
		}(strategy)

	}
}
func (this *StrategyEngine) StartAllStrategies() {
	for _, strategy := range this.Strategies {
		go func(strategy AbstractStrategy) {
			if !strategy.IsInited() {
				log.Infof("[%s] should be init first", strategy.Name())
				return
			}
			if strategy.IsTrading() {
				log.Infof("[%s] has been trading", strategy.Name())
				return
			}

			strategy.OnStart()

		}(strategy)

	}
}

func (this *StrategyEngine) StopAllStrategies() {
	for _, strategy := range this.Strategies {
		go func(strategy AbstractStrategy) {
			strategy.OnStop()
		}(strategy)
	}
}
func (this *StrategyEngine) ProcessTickerData(tickerData bd.TickerData) {
	if len(this.SymbolStrategyMap[tickerData.Symbol]) > 0 {
		for _, strategy := range this.SymbolStrategyMap[tickerData.Symbol] {
			strategy.OnTicker(tickerData)
		}
	}
}

func (this *StrategyEngine) ProcessBarData(barData bd.BarData) {
	log.Info("process bar data")
}

func (this *StrategyEngine) ProcessOrdersData(orderData bd.OrderData) {
	for _, strategy := range this.SymbolStrategyMap[orderData.Symbol] {
		strategy.OnOrders(orderData)
	}
}

func (this *StrategyEngine) ProcessTradeData(tradeData bd.TradeData) {
	for _, strategy := range this.SymbolStrategyMap[tradeData.Symbol] {
		strategy.OnTrade(tradeData)
	}
}

func (this *StrategyEngine) ProcessPositionData(positionData bd.PositionData) {
	symbol := positionData.Symbol
	if len(this.SymbolStrategyMap[positionData.Symbol]) > 0 {
		for _, strategy := range this.SymbolStrategyMap[symbol] {
			strategy.OnPosition(positionData)
		}
	}
}

func (this *StrategyEngine) ProcessBalAndPosData(balAndPosData bd.BalAndPosData) {
	for _, strategy := range this.Strategies {
		strategy.OnBalAndPos(balAndPosData)
	}
}

func (this *StrategyEngine) ProcessBookData(bookData bd.BookData) {
	log.Info("process book data: %v", bookData)
}
