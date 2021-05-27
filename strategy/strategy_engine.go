package strategy

import "sync"

type TickerData struct{}

type StrategyEngine struct {
	strategies []AbstractStrategy
}

var once sync.Once
var instance *StrategyEngine

func GetInstance() *StrategyEngine {
	once.Do(func() {
		instance = NewStrategyEngine()
	})
	return instance
}

func NewStrategyEngine() *StrategyEngine {
	return &StrategyEngine{strategies: make([]AbstractStrategy, 0)}
}

func (this *StrategyEngine) Init()               {}
func (this *StrategyEngine) InitAllStrategies()  {}
func (this *StrategyEngine) StartAllStrategies() {}
func (this *StrategyEngine) ProcessTickerData(tickerData TickerData) {
	println("process ticker data")
}
