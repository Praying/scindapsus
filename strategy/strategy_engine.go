package strategy

type StrategyEngine struct {
	strategies []AbstractStrategy
}

func NewStrategyEngine() *StrategyEngine {
	return &StrategyEngine{strategies: make([]AbstractStrategy, 0)}
}

func (this *StrategyEngine) Init()               {}
func (this *StrategyEngine) InitAllStrategies()  {}
func (this *StrategyEngine) StartAllStrategies() {}
