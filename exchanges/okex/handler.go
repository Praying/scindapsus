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
	case "tickers":
		tickerData := parseTickerData(data)
		event.GetEventEngine().TickerChan <- *tickerData
		return nil
	case "books":
		fallthrough
	case "books5":
		fallthrough
	case "books-l2-tbt":
		fallthrough
	case "books50-l2-tbt":
		bookData := parseBookData(data)
		event.GetEventEngine().BookChan <- *bookData
		return nil
	case "position":
		positionData := parsePositionData(data)
		event.GetEventEngine().PositionChan <- *positionData
		return nil
	case "order":
		orderData := parseOrderData(data)
		event.GetEventEngine().OrderChan <- *orderData
		return nil
	case "orders":
		orderData, tradeData := parseOrdersInfo(data)
		event.GetEventEngine().OrderChan <- *orderData
		event.GetEventEngine().TradeChan <- *tradeData
		return nil
	case "balance_and_position":
		log.Info(data)
		balAndPosData := parseBalAndPosData(data)
		event.GetEventEngine().BalAndPosChan <- *balAndPosData
		return nil
	default:
		return nil
	}
	return nil
}

func (wsClient *OKExWSClient) handle(msg []byte) error {
	//TODO
	log.Info("[ws][response]", string(msg))
	if string(msg) == "pong" {
		return nil
	} else if strings.Contains(string(msg), "clOrdId") {
		//处理下单的状态
		var orderResp OrderResp
		err := json.Unmarshal(msg, &orderResp)
		if err != nil {
			log.Error(err)
			return err
		}
		if orderResp.Code == "1" {
			if len(orderResp.Data) > 0 {
				log.Errorf("clOrdId:%s, ordId:%s,sCode:%s, sMsg:%s", orderResp.Data[0].ClOrdID, orderResp.Data[0].OrdID, orderResp.Data[0].SCode, orderResp.Data[0].SMsg)
				return fmt.Errorf("clOrdId:%s, ordId:%s,sCode:%s, sMsg:%s", orderResp.Data[0].ClOrdID, orderResp.Data[0].OrdID, orderResp.Data[0].SCode, orderResp.Data[0].SMsg)
			}
		} else if orderResp.Code == "0" {
			log.Infof("order successful: ordId: %s", orderResp.Data[0].OrdID)
			return nil
		}

	} else if strings.Contains(string(msg), "event") {
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
			case "error":
				log.Errorf(string(msg))
			default:
				log.Info(string(msg))
			}
			return fmt.Errorf("unknown websocket message: %v", wsResp)
		}
		if wsResp.Code != "" && wsResp.Code != "0" {
			log.Errorf("error, %v", string(msg))

		}

		if wsResp.Arg.Channel != "" {
			return wsClient.respHandler(wsResp.Arg.Channel, msg)
		}

	} else if strings.Contains(string(msg), "arg") {
		//处理推送的
		var wsResp wsResp
		err := json.Unmarshal(msg, &wsResp)
		if err != nil {
			log.Error(err)
			return err
		}

		return wsClient.respHandler(wsResp.Arg.Channel, msg)
	}

	return nil
}
