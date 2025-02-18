package binance

import "github.com/Antkky/go_crypto_scraper/structs"

// outputs
var TestCasesR1 = []structs.TickerData{
	//unwrapped valid message1
	{
		TimeStamp: 1028413123,
		Date:      1823123,
		Symbol:    "BTCUSDT",
		BidPrice:  97245.24,
		BidSize:   53,
		AskPrice:  97260.20,
		AskSize:   42,
	},

	//wrapped valid message1
	{
		TimeStamp: 1028413123,
		Date:      1823123,
		Symbol:    "BTCUSDT",
		BidPrice:  97245.24,
		BidSize:   53,
		AskPrice:  97260.20,
		AskSize:   42,
	},

	//unwrapped invalid message1
	{
		TimeStamp: 1028413123,
		Date:      1823123,
		Symbol:    "BTCUSDT",
		BidPrice:  97245.24,
		BidSize:   53,
		AskPrice:  97260.20,
		AskSize:   42,
	},

	//wrapped invalid message1
	{
		TimeStamp: 1028413123,
		Date:      1823123,
		Symbol:    "BTCUSDT",
		BidPrice:  97245.24,
		BidSize:   53,
		AskPrice:  97260.20,
		AskSize:   42,
	},

	//invalid json message1
	{
		TimeStamp: 1028413123,
		Date:      1823123,
		Symbol:    "BTCUSDT",
		BidPrice:  97245.24,
		BidSize:   53,
		AskPrice:  97260.20,
		AskSize:   42,
	},
}

var TestCasesR2 = []structs.TradeData{
	//unwrapped valid message1
	{},
	//wrapped valid message1
	{},
	//unwrapped invalid message1
	{},
	//wrapped invalid message1
	{},
	//invalid json message1
	{},
}

// inputs
var TestCasesByteArrays = [][]byte{
	//UnwrappedValidMessage
	[]byte(`{
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

	//WrappedValidMessage
	[]byte(`{
		"stream": "BTCUSDT@ticker",
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

	//UnwrappedInvalidMessage
	[]byte(`{
		"stream": "BTCUSDT@ticker",
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

	//WrappedInvalidMessage
	[]byte(`{
		"stream": "BTCUSDT@ticker",
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

	//InvalidJsonMessage
	[]byte(`{invalid json}`),
}
