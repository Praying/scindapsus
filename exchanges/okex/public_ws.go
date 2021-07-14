package okex

type PublicWSClient = OKExWSClient

//订阅行情Ticker数据
func (wsClient *PublicWSClient) WatchTicker(currencyPairs []string) error {
	return wsClient.doTicker(WS_SUBSCRIBE, currencyPairs)
}

func (wsClient *PublicWSClient) UnWatchTicker(currencyPairs []string) error {
	return wsClient.doTicker(WS_UNSUBSCRIBE, currencyPairs)
}

//订阅深度
func (wsClient *PublicWSClient) WatchDepth(currencyPairs []string) error {
	return wsClient.doDepth(WS_SUBSCRIBE, currencyPairs)
}

func (wsClient *PublicWSClient) UnWatchDepth(currencyPairs []string) error {
	return wsClient.doDepth(WS_UNSUBSCRIBE, currencyPairs)
}

//订阅交易
func (wsClient *PublicWSClient) WatchTrades(symbol string) error {
	return wsClient.doTrades(WS_SUBSCRIBE, symbol)
}

func (wsClient *PublicWSClient) UnWatchTrades(symbol string) error {
	return wsClient.doTrades(WS_UNSUBSCRIBE, symbol)
}

//订阅资金费率
func (wsClient *PublicWSClient) WatchFundingRate(symbol string) error {
	return wsClient.doFundingRate(WS_SUBSCRIBE, symbol)
}

func (wsClient *PublicWSClient) UnWatchFundingRate(symbol string) error {
	return wsClient.doFundingRate(WS_UNSUBSCRIBE, symbol)
}

func (wsClient *PublicWSClient) doTrades(op string, symbol string) error {
	var tradesParam TradesParam
	tradesParam.Op = op
	tradesParam.Args = append(tradesParam.Args, struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	}{Channel: TRADES_CHANNEL, InstID: symbol})
	wsClient.WSConn.Subscribe(tradesParam)
	return nil
}

func (wsClient *PublicWSClient) doTicker(op string, currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = op
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{FUNDING_RATE_CHANNEL, currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

//订阅深度数据
func (wsClient *PublicWSClient) doDepth(op string, currencyPairs []string) error {
	var subParam SubParam
	subParam.Op = op
	for _, currencyPair := range currencyPairs {
		subParam.Args = append(subParam.Args, struct {
			Channel string `json:"channel"`
			InstID  string `json:"instId"`
		}{BOOKS5_CHANNEL, currencyPair})
	}
	wsClient.WSConn.Subscribe(subParam)
	return nil
}

func (wsClient *PublicWSClient) doFundingRate(op string, symbol string) error {
	var fundingRateParam FundingRateParam
	fundingRateParam.Op = op
	fundingRateParam.Args = append(fundingRateParam.Args, struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	}{Channel: FUNDING_RATE_CHANNEL, InstID: symbol})
	wsClient.WSConn.Subscribe(fundingRateParam)
	return nil
}
