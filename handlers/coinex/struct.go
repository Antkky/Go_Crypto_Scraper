package coinex

import "encoding/json"

type GlobalMessageStruct struct {
	Method  string          `json:"method"`
	Data    json.RawMessage `json:"data"`
	Id      int             `json:"id"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
}

type TickerData struct {
	Method  string            `json:"method"`
	Data    TickerDataPayload `json:"data"`
	Id      int               `json:"id"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
}

type TradeData struct {
	Method  string           `json:"method"`
	Data    TradeDataPayload `json:"data"`
	Id      int              `json:"id"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
}

type TickerDataPayload struct {
	Market     string `json:"market"`
	Updated_at int    `json:"updated_at"`
	BidPrice   string `json:"best_bid_price"`
	BidSize    string `json:"best_bid_size"`
	AskPrice   string `json:"best_ask_price"`
	AskSize    string `json:"best_ask_size"`
}

type TradeDataPayload struct {
	Market string  `json:"market"`
	Deals  []Trade `json:"deal_list"`
}

type Trade struct {
	ID         int    `json:"deal_id"`
	Created_at int    `json:"created_at"`
	Side       string `json:"side"`
	Price      string `json:"price"`
	Amount     string `json:"amount"`
}
