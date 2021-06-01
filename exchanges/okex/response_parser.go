package okex

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"scindapsus/strategy"
	"strconv"
	"time"
)

func stringTof64(input string) float64 {
	res, _ := strconv.ParseFloat(input, 64)
	return res
}
func stringToInt64(input string) int64 {
	res, _ := strconv.ParseInt(input, 10, 64)
	return res
}

//根据交易所返回的Ticker数据返回统一的TickerData格式
func parseTickerData(data []byte) *strategy.TickerData {
	var tickerResp TickerResp
	err := json.Unmarshal(data, &tickerResp)
	if err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to TickerResp failed", err.Error(), string(data))
		return nil
	}
	tickerData := &strategy.TickerData{
		Symbol:        tickerResp.Arg.Instid,
		TimeStamp:     stringToInt64(tickerResp.Data[0].Ts),
		DateTime:      time.Unix(stringToInt64(tickerResp.Data[0].Ts), 0),
		High:          stringTof64(tickerResp.Data[0].High24H),
		Low:           stringTof64(tickerResp.Data[0].Low24H),
		Open:          stringTof64(tickerResp.Data[0].Open24H),
		Last:          stringTof64(tickerResp.Data[0].Last),
		LastVolume:    stringTof64(tickerResp.Data[0].Lastsz),
		Close:         0,
		PreviousClose: 0,
		Bid:           stringTof64(tickerResp.Data[0].Bidpx),
		BidVolume:     stringTof64(tickerResp.Data[0].Bidsz),
		Ask:           stringTof64(tickerResp.Data[0].Askpx),
		AskVolume:     stringTof64(tickerResp.Data[0].Asksz),
		BaseVolume:    stringTof64(tickerResp.Data[0].Vol24H),
		QuoteVolume:   stringTof64(tickerResp.Data[0].Volccy24H),
		Percentage:    0,
		Average:       0,
		Change:        0,
		VWap:          0,
	}
	return tickerData
}
