package binance

import (
	"github.com/Antkky/go_crypto_scraper/utils"
)

// Test Cases for ProcessMessageType
var ProcessMessageTypeCases = []struct {
	name       string
	eventType  string
	message    []byte
	wrapped    bool
	r1         utils.TickerDataStruct
	r2         utils.TradeDataStruct
	errorValue error
	wantError  bool
}{
	// Unwrapped valid message
	{
		name:      "unwrapped valid message1",
		eventType: "24hrTicker",
		message: []byte(`{
			"e": "24hrTicker",
			"E": 1672515782136,
			"s": "BNBBTC",
			"p": "0.0015",
			"P": "250.00",
			"w": "0.0018",
			"x": "0.0009",
			"c": "0.0025",
			"Q": "10",
			"b": "0.0024",
			"B": "10",
			"a": "0.0026",
			"A": "100",
			"o": "0.0010",
			"h": "0.0025",
			"l": "0.0010",
			"v": "10000",
			"q": "18",
			"O": 0,
			"C": 86400000,
			"F": 0,
			"L": 18150,
			"n": 18151
		}`),
		wrapped: false,
		r1: utils.TickerDataStruct{
			TimeStamp: 1672515782136,
			Date:      1672515782136,
			Symbol:    "BNBBTC",
			BidPrice:  "0.0024",
			BidSize:   "10",
			AskPrice:  "0.0026",
			AskSize:   "100",
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  false,
	},
}

// Test Cases for HandleMessage
var HandleMessageTestCases = []struct {
	name       string
	message    []byte
	exchange   utils.ExchangeConfig
	errorValue error
	wantError  bool
}{
	{
		name: "valid Global Message #1",
		message: []byte(`{
			"stream": "btcusdt@ticker",
			"data": {
				"e": "24hrTicker",
				"E": 1672515782136,
				"s": "BNBBTC",
				"p": "0.0015",
				"P": "250.00",
				"w": "0.0018",
				"x": "0.0009",
				"c": "0.0025",
				"Q": "10",
				"b": "0.0024",
				"B": "10",
				"a": "0.0026",
				"A": "100",
				"o": "0.0010",
				"h": "0.0025",
				"l": "0.0010",
				"v": "10000",
				"q": "18",
				"O": 0,
				"C": 86400000,
				"F": 0,
				"L": 18150,
				"n": 18151
			}
		}`),
		exchange: utils.ExchangeConfig{
			Name: "Binance Global",
			URI:  "wss://data-stream.binance.vision/stream",
		},
		errorValue: nil,
		wantError:  false,
	},
}
