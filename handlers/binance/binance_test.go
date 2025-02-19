package binance

import (
	"testing"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/stretchr/testify/assert"
)

// _____________Table of Contents______________
// var  ProcessMessageTypeTestCases : line-019
// func TestProcessMessageType()    : line-236
// var  HandleMessageTestCases      : line-269
// func TestHandleMessage()         : line-289
// func TestHandleTickerMessage()   : line-317
// func TestHandleTradeMessage()    : line-364


// Test Cases for ProcessMessageType
var ProcessMessageTypeTestCases = []struct {
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
		eventType: "24hrMiniTicker",
		message: []byte(`{}
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
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  false,
	},

	// Wrapped valid message
	{
		name:      "wrapped valid message1",
		eventType: "24hrMiniTicker",
		message: []byte(`{
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
		wrapped: true,
		r1: utils.TickerDataStruct{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  false,
	},

	// Unwrapped invalid message
	{
		name:      "unwrapped invalid message1",
		eventType: "24hrMiniTicker",
		message: []byte(`{
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
		wrapped: false,
		r1: utils.TickerDataStruct{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  true,
	},

	// Wrapped invalid message
	{
		name:      "wrapped invalid message1",
		eventType: "24hrMiniTicker",
		message: []byte(`{
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
		wrapped: true,
		r1: utils.TickerDataStruct{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  true,
	},

	// Invalid JSON message
	{
		name:      "invalid json message1",
		eventType: "24hrMiniTicker",
		message:   []byte(`{invalid json}`),
		wrapped:   false,
		r1: utils.TickerDataStruct{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         utils.TradeDataStruct{},
		errorValue: nil,
		wantError:  true,
	},
}

//TestProcessMessageType
//
// inputs
// message : []byte
// r1      : *utils.TickerDataStruct
// r2      : *utils.TradeDataStruct
//
// Outputs:
// err : error
//
// Description:
// routes message for processing changes through pointer reference
func TestProcessMessageType(t *testing.T) {
	for _, tt := range ProcessMessageTypeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				r1 utils.TickerDataStruct
				r2 utils.TradeDataStruct
			)

			err := ProcessMessageType(tt.message, &r1, &r2)

			// error test
			if assert.Error(t, err, "An error has occurred") {
				if !tt.wantError && assert.Equal(t, tt.errorValue, err, "Unexpected error") {
					t.Errorf("Unexpected error: %+v", err)
				}
			}

			// r1 test
			if !assert.Equal(t, r1, tt.r1) && !tt.wantError {
				t.Errorf("r1 isn't expected\nr1: %+v\n expected: %+v", r1, tt.r1)
			}

			// r2 test
			if !assert.Equal(t, r2, tt.r2) && !tt.wantError {
				t.Errorf("r2 isn't expected\nr2: %+v\n expected: %+v", r2, tt.r2)
			}
		})
	}
}

// Test Cases for HandleMessage
var HandleMessageTestCases = []struct {
	name       string
	message    []byte
	exchange   utils.ExchangeConfig
	errorValue error
	wantError  bool
}{
	{},
}

//TestHandleMessage
//
// inputs
// message : []byte
// exchange      : *utils.ExchangeConfig
//
// Outputs:
// err : error
//
// Description:
// routes message for processing
// based on type of message
func TestHandleMessage(t *testing.T) {
	for _, tt := range HandleMessageTestCases {
		t.Run(tt.name, func(t *testing.T) {
			// Run Test
			err := HandleMessage(tt.message, tt.exchange)

			// Test Errors
			if assert.Error(t, err, "An error has occurred") {
				if !tt.wantError && !assert.Equal(t, tt.errorValue, err, "Expected Error") {
					t.Errorf("Unexpected Error: %+v", err)
				}
			}
		})
	}
}
