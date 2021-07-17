package okex

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"scindapsus/event"
	"strings"
)

func OkexRespHandler(channel string, data json.RawMessage) error {
	switch channel {
	case TICKER_CHANNEL:
		tickerData := parseTickerData(data)
		event.GetEventEngine().TickerChan <- *tickerData
		return nil
	case BOOKS_CHANNEL:
		fallthrough
	case BOOKS5_CHANNEL:
		fallthrough
	case BOOKS_L2_TBT_CHANNEL:
		fallthrough
	case BOOKS50_L2_TBT_CHANNEL:
		bookData := parseBookData(data)
		event.GetEventEngine().BookChan <- *bookData
		return nil
	case POSITIONS_CHANNEL:
		positionData := parsePositionData(data)
		if positionData != nil {
			event.GetEventEngine().PositionChan <- *positionData
		}
		return nil
	case ORDER_CHANNEL:
		//orderData := parseOrderData(data)
		//event.GetEventEngine().OrdersChan <- *orderData
		return nil
	case ORDERS_CHANNEL:
		orderData, tradeData := parseOrdersInfo(data)
		if orderData != nil {
			event.GetEventEngine().OrdersChan <- *orderData
		}
		if tradeData != nil {
			event.GetEventEngine().TradeChan <- *tradeData
		}
		return nil
	case BAL_AND_POS_CHANNEL:
		log.Info(data)
		balAndPosData := parseBalAndPosData(data)
		event.GetEventEngine().BalAndPosChan <- *balAndPosData
		return nil
	case FUNDING_RATE_CHANNEL:
		fundingRateData := parseFundingRate(data)
		if fundingRateData != nil {
			log.Infof("Funding rate: %f, next funding rate: %f, time: %s", fundingRateData.FundingRate, fundingRateData.NextFundingRate, fundingRateData.FundingTime.String())
		}
		return nil
	case INSTRUMENTS_CHANNEL:
		//TODO 解析推送的数据
		parseInstrument(data)
		return nil
	default:
		return nil
	}
	return nil
}

func (wsClient *OKExWSClient) handle(msg []byte) error {
	msgStr := string(msg)
	//log.Info("[ws][response]", string(msg))
	if msgStr == "pong" {
		return nil
	}

	if strings.Contains(msgStr, "event") {
		//处理订阅事件的状态
		var wsResp wsResp
		err := json.Unmarshal(msg, &wsResp)
		if err != nil {
			log.Error(err)
			return err
		}

		if wsResp.Event != "" {
			switch wsResp.Event {
			case "subscribe":
				log.Info("subscribed:", wsResp.Arg.Channel)
				return nil
			case "unsubscribe":
				log.Info("unsubscribed:", wsResp.Arg.Channel)
				return nil
			case "login":
				log.Info("login:", string(msg))
				return nil
			case "error":
				log.Errorf(string(msg))
				return nil
			default:
				log.Info(string(msg))
			}
			return fmt.Errorf("unknown websocket message: %v", wsResp)

		}
		if wsResp.Code != "" && wsResp.Code != "0" {
			log.Errorf("error, %v", string(msg))

		}
	}

	if strings.Contains(msgStr, "arg") && strings.Contains(msgStr, "channel") {
		//推送数据
		var wsResp wsResp
		err := json.Unmarshal(msg, &wsResp)
		if err != nil {
			log.Error(err)
			return err
		}
		return wsClient.respHandler(wsResp.Arg.Channel, msg)
	}
	//处理订单操作
	if strings.Contains(msgStr, "op") && strings.Contains(msgStr, "order") {
		var orderResp OrderResp
		if err := json.Unmarshal(msg, &orderResp); err != nil {
			log.Errorf("[handle] failed to  unmarshal order response")
			return nil
		}
		if orderResp.Code == "0" {
			//下单成功
		} else {
			//下单失败，属于正常信息
			log.Infof("order failed, %+v", orderResp)
		}
	}

	return nil
}
