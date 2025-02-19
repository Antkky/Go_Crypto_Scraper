package binance

import (
	"testing"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/stretchr/testify/assert"
)

// _____________Main Stuff_____________

// ProcessMessageType
var ProcessMessageTypeTestCases = []struct {
	name       string
	eventType  string
	message    []byte
	wrapped    bool
	r1         structs.TickerData
	r2         structs.TradeData
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
		r1: structs.TickerData{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         structs.TradeData{},
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
		r1: structs.TickerData{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         structs.TradeData{},
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
		r1: structs.TickerData{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         structs.TradeData{},
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
		r1: structs.TickerData{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         structs.TradeData{},
		errorValue: nil,
		wantError:  true,
	},

	// Invalid JSON message
	{
		name:      "invalid json message1",
		eventType: "24hrMiniTicker",
		message:   []byte(`{invalid json}`),
		wrapped:   false,
		r1: structs.TickerData{
			TimeStamp: 1028413123,
			Date:      1823123,
			Symbol:    "BTCUSDT",
			BidPrice:  97245.24,
			BidSize:   53,
			AskPrice:  97260.20,
			AskSize:   42,
		},
		r2:         structs.TradeData{},
		errorValue: nil,
		wantError:  true,
	},
}

func TestProcessMessageType(t *testing.T) {
	for _, tt := range ProcessMessageTypeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				r1 structs.TickerData
				r2 structs.TradeData
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

// HandleMessage
var HandleMessageTestCases = []struct {
	name       string
	message    []byte
	exchange   structs.ExchangeConfig
	errorValue error
	wantError  bool
}{
	{},
}

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

// _____________Ticker & Trade Handlers_____________
func TestHandleTickerMessage(t *testing.T) {
	tests := []struct {
		name       string
		message    []byte
		wrapped    bool
		r1         structs.TickerData
		errorValue error
		wantError  bool
	}{
		// Test Cases
		{},
	}

	for _, tt := range tests {
		var r1 structs.TickerData
		err := HandleTickerMessage(tt.message, &r1)

		// error test
		if tt.wantError {
			if !assert.Error(t, err) && !assert.Equal(t, tt.errorValue, err) {
				t.Error("Unexpected error")
			}
		} else {
			if assert.Error(t, err) {
				t.Error("Unexpected error")
			}
		}

		if !assert.Equal(t, tt.r1, r1) && !tt.wantError {
			t.Errorf("Unexpected R1\nr1: %+v\nexpected: %+v", r1, tt.r1)
		}
	}
}
func TestHandleTradeMessage(t *testing.T) {
	tests := []struct {
		name       string
		message    []byte
		wrapped    bool
		r1         structs.TradeData
		errorValue error
		wantError  bool
	}{
		// Test Cases
		{},
	}

	for _, tt := range tests {
		var r1 structs.TradeData
		err := HandleTradeMessage(tt.message, &r1)

		// error test
		if tt.wantError {
			if !assert.Error(t, err) && !assert.Equal(t, tt.errorValue, err) {
				t.Error("Unexpected error")
			}
		} else {
			if assert.Error(t, err) {
				t.Error("Unexpected error")
			}
		}

		if !assert.Equal(t, tt.r1, r1) && !tt.wantError {
			t.Errorf("Unexpected R1\nr1: %+v\nexpected: %+v", r1, tt.r1)
		}
	}
}
