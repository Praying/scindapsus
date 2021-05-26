package strategy

type AbstractStrategy interface {
	OnTicker()
	OnTrade()
	OnOrder()
	OnInit()
	OnBar()
	OnStart()
}
