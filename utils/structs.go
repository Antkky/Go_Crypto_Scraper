package utils

type ExchangeConfig struct {
	Name    string                   `json:"name"`
	URI     string                   `json:"uri"`
	Streams []map[string]interface{} `json:"streams"`
	Ping    map[string]interface{}   `json:"ping,omitempty"`
}

// add some functionallity
type TickerDataStruct struct {
	TimeStamp uint
	Date      uint
	Symbol    string
	BidPrice  float32
	BidSize   float32
	AskPrice  float32
	AskSize   float32
}

type TradeDataStruct struct {
	TimeStamp uint
	Date      uint
	Symbol    string
	Price     float32
	Quantity  float32
	Bid_MM    bool
}

// Buffer Structs
type DataBuffer struct {
	tickerBuffer [][]TickerDataStruct
	tradeBuffer  [][]TradeDataStruct
	maxSize      int
	dataType     string
	dataStream   string
	filePath     string
	fileName     string
}
