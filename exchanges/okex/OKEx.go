package okex

type Exchange interface {
	SendLimitOrder()
}

type OKEx struct {
	//Rest和WebSocket

}
