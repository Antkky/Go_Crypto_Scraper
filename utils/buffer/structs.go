package buffer

import "github.com/Antkky/go_crypto_scraper/utils"

// Buffer Structs
type DataBuffer struct {
	TickerBuffer []utils.TickerDataStruct
	TradeBuffer  []utils.TradeDataStruct
	MaxSize      int
	DataType     string
	Symbol       string
	Market       string
	ID           string
	FilePath     string
	FileName     string
}
