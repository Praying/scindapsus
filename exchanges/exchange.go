package exchanges

type Exchange interface {
	Buy()
	Sell()
	Long()
	Short()
	Cover()
	Init()
}

type OKExchange struct {
}

func NewOKExchange() *OKExchange {
	return &OKExchange{}
}

func (O OKExchange) Init() {
	panic("implement me")
}

func (O OKExchange) Buy() {
	panic("implement me")
}

func (O OKExchange) Sell() {
	panic("implement me")
}

func (O OKExchange) Long() {
	panic("implement me")
}

func (O OKExchange) Short() {
	panic("implement me")
}

func (O OKExchange) Cover() {
	panic("implement me")
}
