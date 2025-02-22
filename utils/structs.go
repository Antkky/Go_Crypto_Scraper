package utils

import (
	"encoding/json"
)

type ExchangeConfig struct {
	Name    string `json:"name"`
	URI     string `json:"uri"`
	Market  string `json:"market"`
	Streams []struct {
		Type    string          `json:"type"`
		Symbol  string          `json:"symbol"`
		Market  string          `json:"market"`
		Message json.RawMessage `json:"message"`
	} `json:"streams"`
	Ping map[string]interface{} `json:"ping,omitempty"`
}

// add some functionallity
type TickerDataStruct struct {
	TimeStamp uint64
	Date      uint64
	Symbol    string
	BidPrice  string
	BidSize   string
	AskPrice  string
	AskSize   string
}

type TradeDataStruct struct {
	TimeStamp uint64
	Date      uint64
	Symbol    string
	Price     string
	Quantity  string
	Bid_MM    bool
}

// Buffer Structs
type DataBuffer struct {
	tickerBuffer [][]TickerDataStruct
	tradeBuffer  [][]TradeDataStruct
	maxSize      uint16
	dataType     string
	dataStream   string
	//filePath     string
	fileName string
}
