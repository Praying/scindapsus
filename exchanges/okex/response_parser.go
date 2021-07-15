package okex

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	bd "scindapsus/basedata"

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
func parseTickerData(data []byte) *bd.TickerData {
	var tickerResp TickerResp
	err := json.Unmarshal(data, &tickerResp)
	if err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to TickerResp failed", err.Error(), string(data))
		return nil
	}
	tickerData := &bd.TickerData{
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

func parseBookData(data []byte) *bd.BookData {
	var bookResp BookResp
	if err := json.Unmarshal(data, &bookResp); err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to BookResp failed", err.Error(), string(data))
		return nil
	}

	bookData := &bd.BookData{
		Symbol:  bookResp.Arg.InstID,
		Action:  bookResp.Action,
		AskList: make([]bd.BookRecord, 0),
		BidList: make([]bd.BookRecord, 0),
	}
	if len(bookResp.Data) == 1 {
		for _, ask := range bookResp.Data[0].Asks {
			bookRecord := bd.BookRecord{
				Price:  stringTof64(ask[0]),
				Volume: stringTof64(ask[1]),
			}
			bookData.AskList = append(bookData.AskList, bookRecord)
		}
		for _, bid := range bookResp.Data[0].Bids {
			bookRecord := bd.BookRecord{
				Price:  stringTof64(bid[0]),
				Volume: stringTof64(bid[1]),
			}
			bookData.BidList = append(bookData.BidList, bookRecord)
		}
	}

	return bookData
}

func parsePositionData(data []byte) *bd.PositionData {
	var positionResp PositionResp
	if err := json.Unmarshal(data, &positionResp); err != nil {
		log.Errorf("error:%s \nUnmarshal data:%s to PositionResp failed", err.Error(), string(data))
		return nil
	}
	for _, posData := range positionResp.Data {
		positionData := &bd.PositionData{
			Pos: stringTof64(posData.Pos),
		}
		return positionData
	}
	return nil
}

func parseOrderData(data []byte) *bd.OrderData {
	orderData := &bd.OrderData{}
	return orderData
}

func parseBalAndPosData(data []byte) *bd.BalAndPosData {
	var balAndPosResp BalAndPosResp
	if err := json.Unmarshal(data, &balAndPosResp); err != nil {
		log.Errorf("error: %s, Unmarshal data: %s to BalAndPostionResp", err.Error(), data)
		return nil
	}
	balAndPosData := &bd.BalAndPosData{
		BalMap:      make(map[string]float64),
		PositionMap: make(map[string]float64),
	}

	if len(balAndPosResp.Data) > 0 {
		for _, data := range balAndPosResp.Data {
			if len(data.BalData) > 0 {
				for _, balData := range data.BalData {
					balAndPosData.BalMap[balData.Ccy] = stringTof64(balData.CashBal)
				}
			}

			if len(data.PosData) > 0 {
				for _, posData := range data.PosData {
					balAndPosData.PositionMap[posData.Ccy] = stringTof64(posData.Pos)
				}
			}
		}
	}
	return balAndPosData
	return nil
}

/*
   def on_order(self, packet: dict) -> None:
       """委托更新推送"""
       data = packet["data"]
       for d in data:
           order: OrderData = parse_order_data(d, self.gateway_name)
           self.gateway.on_order(order)

           # 检查是否有成交
           if d["fillSz"] == "0":
               return

           # 将成交数量四舍五入到正确精度
           trade_volume: float = float(d["fillSz"])
           contract: ContractData = symbol_contract_map.get(order.symbol, None)
           if contract:
               trade_volume = round_to(trade_volume, contract.min_volume)

           trade: TradeData = TradeData(
               symbol=order.symbol,
               exchange=order.exchange,
               orderid=order.orderid,
               tradeid=d["tradeId"],
               direction=order.direction,
               offset=order.offset,
               price=float(d["fillPx"]),
               volume=trade_volume,
               datetime=parse_timestamp(d["uTime"]),
               gateway_name=self.gateway_name,
           )
           self.gateway.on_trade(trade)
*/
func parseOrdersInfo(data []byte) (*bd.OrderData, *bd.TradeData) {
	var ordersInfo OrdersInfo
	if err := json.Unmarshal(data, &ordersInfo); err != nil {
		log.Errorf("error: %s, Unmarshal data: %s to BalAndPostionResp", err.Error(), data)
		return nil, nil
	}
	for _, item := range ordersInfo.Data {
		//检查是否哟成交
		orderId := item.ClOrdID
		if orderId == "" {
			orderId = item.OrdID
		}
		orderData := &bd.OrderData{
			Symbol:    item.InstID,
			Exchange:  bd.Exchange_OKEX,
			OrderID:   orderId,
			OrderType: bd.OrderType_Okex2Vt[item.OrdType],
			Direction: bd.Direction_Okex2Vt[item.Side],
			Offset:    "",
			Price:     stringTof64(item.Px),
			Volume:    stringTof64(item.Sz),
			Traded:    stringTof64(item.AccFillSz),
			Status:    bd.Status_Okex2Vt[item.State],
			DateTime:  parseTime(item.CTime),
			Reference: "",
		}
		if item.FillSz == "0" {
			return orderData, nil
		}
		tradeData := &bd.TradeData{
			Symbol:    item.InstID,
			Exchange:  bd.Exchange_OKEX,
			OrderID:   orderId,
			TradeID:   item.TradeID,
			Direction: bd.Direction_Okex2Vt[item.Side],
			Price:     stringTof64(item.FillPx),
			Volume:    stringTof64(item.FillSz),
			Datetime:  parseTime(item.CTime),
			STSymbol:  "",
			STOrderID: "",
			STTradeID: "",
		}
		return orderData, tradeData
	}
	return nil, nil
}

func parseFundingRate(data []byte) *bd.FundingRateData {
	var fundingRateInfo FundingRateInfo
	if err := json.Unmarshal(data, &fundingRateInfo); err != nil {
		log.Errorf("error: %s, Unmarshal data: %s to FundingRateInfo failed", err.Error(), data)
		return nil
	}
	if len(fundingRateInfo.Data) > 0 {
		for _, item := range fundingRateInfo.Data {
			ts, _ := strconv.ParseInt(item.FundingTime, 10, 64)
			fundingRateData := &bd.FundingRateData{
				FundingRate:     stringTof64(item.FundingRate),
				NextFundingRate: stringTof64(item.NextFundingRate),
				FundingTime:     time.Unix(ts/1000, (ts%1000)*1000000),
			}
			return fundingRateData
		}
	}
	return nil
}

//@param 传入的timestamp是13位的时间戳
func parseTime(timestamp string) time.Time {
	//s:="1597026383085"
	data, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Errorf("parse %v failed", timestamp)
		return time.Time{}
	}
	msec := data % 1000
	return time.Unix(data/1000, msec*1000000)
}

func parseInstrument(data []byte) error {
	var instrumentsResp InstrumentsResp
	if err := json.Unmarshal(data, &instrumentsResp); err != nil {
		log.Errorf("error: %s, Unmarshal data: %s to InstrumentsResp failed", err.Error(), data)
		return nil
	}
	if instrumentsResp.Arg.InstType == INST_SWAP {
		for _, item := range instrumentsResp.Data {
			if item.InstID == "ETH-USD-SWAP" {
				log.Infof("%+v\n", item)
			}
		}
	}
	return nil
}
