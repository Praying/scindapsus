package okex

//推送数据参数
type TickerResp struct {
	Arg struct {
		Channel string `json:"channel"`
		Instid  string `json:"instId"`
	} `json:"arg"`
	Data []struct {
		Insttype  string `json:"instType"`
		Instid    string `json:"instId"`
		Last      string `json:"last"`
		Lastsz    string `json:"lastSz"`
		Askpx     string `json:"askPx"`
		Asksz     string `json:"askSz"`
		Bidpx     string `json:"bidPx"`
		Bidsz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		Volccy24H string `json:"volCcy24h"`
		Vol24H    string `json:"vol24h"`
		Sodutc0   string `json:"sodUtc0"`
		Sodutc8   string `json:"sodUtc8"`
		Ts        string `json:"ts"`
	} `json:"data"`
}
