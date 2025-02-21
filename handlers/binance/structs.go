package binance

import "encoding/json"

// Global Message Struct
type GlobalMessageStruct struct {
	Stream string `json:"stream"`
	Data   struct {
		EventType string `json:"e"`
		EventTime int64  `json:"E"`
		Symbol    string `json:"s"`
	} `json:"data"`
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Result    string `json:"result"`
}

type GlobalMessagePayload struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Result    string `json:"result"`
}

type TickerData struct {
	EventType   string      `json:"e"`
	EventTime   int64       `json:"E"`
	Symbol      string      `json:"s"`
	BidPrice    json.Number `json:"b"`
	BidSize     json.Number `json:"B"`
	AskPrice    json.Number `json:"a"`
	AskSize     json.Number `json:"A"`
	ClosePrice  json.Number `json:"c"`
	OpenPrice   json.Number `json:"o"`
	HighPrice   json.Number `json:"h"`
	LowPrice    json.Number `json:"l"`
	BaseVolume  json.Number `json:"v"`
	QuoteVolume json.Number `json:"q"`
}

type TradeData struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int    `json:"t"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int64  `json:"T"`
	IsMaker   bool   `json:"m"`
	Ignore    bool   `json:"M"`
}
