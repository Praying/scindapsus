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
		Volccy24H string `json:"volCcy24h"` //Ccy是以币为单位
		Vol24H    string `json:"vol24h"`
		Sodutc0   string `json:"sodUtc0"`
		Sodutc8   string `json:"sodUtc8"`
		Ts        string `json:"ts"`
	} `json:"data"`
}

//深度数据
type BookResp struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Action string `json:"action"`
	Data   []struct {
		Asks     [][]string `json:"asks"`
		Bids     [][]string `json:"bids"`
		Ts       string     `json:"ts"`
		Checksum int        `json:"checksum"`
	} `json:"data"`
}

type SubParam struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"args"`
}
