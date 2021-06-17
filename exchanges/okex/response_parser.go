package okex

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"scindapsus/strategy"
	"strconv"
	"time"
)

func stringTof64(input string) float64 {
	res, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.Errorf("parse %s to float64 got error:%s", input, err.Error())
		return 0.0
	}
	return res
}
func stringToInt64(input string) int64 {
	res, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		log.Errorf("parse %s to int64 got error:%s", input, err.Error())
		return 0
	}
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

func parseBookData(data []byte) *strategy.BookData {
	var bookResp BookResp
	if err := json.Unmarshal(data, &bookResp); err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to BookResp failed", err.Error(), string(data))
		return nil
	}

	bookData := &strategy.BookData{
		Symbol:  bookResp.Arg.InstID,
		Action:  bookResp.Action,
		AskList: make([]strategy.BookRecord, 0),
		BidList: make([]strategy.BookRecord, 0),
	}
	if len(bookResp.Data) == 1 {
		for _, ask := range bookResp.Data[0].Asks {
			bookRecord := strategy.BookRecord{
				Price:  stringTof64(ask[0]),
				Volume: stringTof64(ask[1]),
			}
			bookData.AskList = append(bookData.AskList, bookRecord)
		}
		for _, bid := range bookResp.Data[0].Bids {
			bookRecord := strategy.BookRecord{
				Price:  stringTof64(bid[0]),
				Volume: stringTof64(bid[1]),
			}
			bookData.BidList = append(bookData.BidList, bookRecord)
		}
	}

	return bookData
}

func parsePositionData(data []byte) *strategy.PositionData {
	var positionResp PositionResp
	if err := json.Unmarshal(data, &positionResp); err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to PositionResp failed", err.Error(), string(data))
		return nil
	}
	positionData := &strategy.PositionData{}
	return positionData
}

func parseOrderData(data []byte) *strategy.OrderData {
	orderData := &strategy.OrderData{}
	return orderData
}

func parseBalAndPosData(data []byte) interface{} {
	var balAndPosResp BalAndPosResp
	if err := json.Unmarshal(data, &balAndPosResp); err != nil {
		log.Errorf("error: %s, Unmarshal data: %s to BalAndPostionResp", err.Error(), data)
		return nil
	}
	log.Infof("%v", balAndPosResp)
	return nil
}
