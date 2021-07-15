package okex

const (
	//公共产品频道
	INSTRUMENTS_CHANNEL string = "instruments"
	//公共行情频道
	TICKER_CHANNEL string = "tickers"
	//公共交易频道
	TRADES_CHANNEL string = "trades"
	//交易操作
	ORDER_CHANNEL string = "order"
	//订单频道
	ORDERS_CHANNEL string = "orders"
	//资金费率频道
	FUNDING_RATE_CHANNEL string = "funding-rate"
	//持仓频道
	POSITIONS_CHANNEL string = "positions"
	//持仓和余额
	BAL_AND_POS_CHANNEL string = "balance_and_position"

	/*
		books 首次推400档快照数据，以后增量推送，即每100毫秒有深度变化推送一次变化的数据
		books5首次推5档快照数据，以后定量推送，每100毫秒有深度变化推送一次5档数据，即每次都推送5档数据
		books-l2-tbt 首次推400档快照数据，以后增量推送，即有深度有变化推送一次变化的数据
		books50-l2-tbt 首次推50档快照数据，以后增量推送，即有深度有变化推送一次变化的数据
	*/
	BOOKS_CHANNEL          string = "books"
	BOOKS5_CHANNEL         string = "books5"
	BOOKS_L2_TBT_CHANNEL   string = "books-l2-tbt"
	BOOKS50_L2_TBT_CHANNEL string = "books50-l2-tbt"
)
