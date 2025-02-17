package binance

import (
	"testing"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name     string
	message  []byte
	wrapped  bool
	expected structs.TickerData
}

func TestHandleTickerMessage(t *testing.T) {
	UnwrappedValidMessage := TestCase{
		name: "unwrapped valid message",
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
			}
		`),
		wrapped: false,
		expected: structs.TickerData{
			TimeStamp: 1672515782136,
			Date:      1212,
			Symbol:    "BTCUSDT",
			BidPrice:  0.0024,
			BidSize:   10,
			AskPrice:  0.0026,
			AskSize:   100,
		},
	}
	WrappedValidMessage := TestCase{
		name: "unwrapped valid message",
		message: []byte(`
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
		`),
		wrapped: true,
		expected: structs.TickerData{
			TimeStamp: 1672515782136,
			Date:      1212,
			Symbol:    "BTCUSDT",
			BidPrice:  0.0024,
			BidSize:   10,
			AskPrice:  0.0026,
			AskSize:   100,
		},
	}
	UnwrappedInvalidMessage := TestCase{
		name: "unwrapped valid message",
		message: []byte(`
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
		`),
		wrapped: true,
		expected: structs.TickerData{
			TimeStamp: 1672515782136,
			Date:      1212,
			Symbol:    "BTCUSDT",
			BidPrice:  0.0024,
			BidSize:   10,
			AskPrice:  0.0026,
			AskSize:   100,
		},
	}
	WrappedInvalidMessage := TestCase{
		name: "unwrapped valid message",
		message: []byte(`
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
		`),
		wrapped: true,
		expected: structs.TickerData{
			TimeStamp: 1672515782136,
			Date:      1212,
			Symbol:    "BTCUSDT",
			BidPrice:  0.0024,
			BidSize:   10,
			AskPrice:  0.0026,
			AskSize:   100,
		},
	}
	InvalidJsonMessage := TestCase{
		name:     "invalid JSON message",
		message:  []byte(`{invalid json}`),
		wrapped:  false,
		expected: structs.TickerData{},
	}

	tests := []TestCase{
		UnwrappedValidMessage,
		WrappedValidMessage,
		UnwrappedInvalidMessage,
		WrappedInvalidMessage,
		InvalidJsonMessage,
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleTickerMessage(tt.message, tt.wrapped)
			if !assert.Equal(t, result, tt.expected, "Should Be Equal") {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestHandleTradeMessage(t *testing.T) {

}
