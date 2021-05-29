package okex

import "net/http"

type Exchange interface {
	SendLimitOrder()
}

type OKExWSClient struct {
}

type OKExRestClient struct {
}

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro

	Lever float64 //杠杆倍数 , for future
}
type OKExExchange struct {
	//Rest和WebSocket
	publicWS   *OKExWSClient
	privateWS  *OKExWSClient
	restClient *OKExRestClient
	config     *APIConfig
}

func NewOKExExchange(config *APIConfig) *OKExExchange {
	return &OKExExchange{config: config}
}
