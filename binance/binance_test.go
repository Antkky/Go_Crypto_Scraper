package binance

import (
	"testing"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/stretchr/testify/assert"
)

// Main Stuff
func TestProcessMessageType(t *testing.T) {
	tests := []struct {
		name       string
		eventType  string
		message    []byte
		wrapped    bool
		r1         structs.TickerData
		r2         structs.TradeData
		errorValue error
		wantError  bool
	}{
		// Test Cases
		{name: "unwrapped valid message1", eventType: "24hrMiniTicker", message: TestCasesByteArrays[0], wrapped: false, r1: TestCasesR1[0], r2: TestCasesR2[0], errorValue: nil, wantError: false},
		{name: "wrapped valid message1", eventType: "24hrMiniTicker", message: TestCasesByteArrays[1], wrapped: true, r1: TestCasesR1[1], r2: TestCasesR2[1], errorValue: nil, wantError: false},
		{name: "unwrapped invalid message1", eventType: "24hrMiniTicker", message: TestCasesByteArrays[2], wrapped: false, r1: TestCasesR1[2], r2: TestCasesR2[2], errorValue: nil, wantError: true},
		{name: "wrapped invalid message1", eventType: "24hrMiniTicker", message: TestCasesByteArrays[3], wrapped: true, r1: TestCasesR1[3], r2: TestCasesR2[3], errorValue: nil, wantError: true},
		{name: "invalid json message1", eventType: "24hrMiniTicker", message: TestCasesByteArrays[4], wrapped: false, r1: TestCasesR1[4], r2: TestCasesR2[4], errorValue: nil, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, r2, err := ProcessMessageType(tt.eventType, tt.message, tt.wrapped)

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
func TestHandleMessage(t *testing.T) {
	// declare tests
	tests := []struct {
		name       string
		message    []byte
		exchange   structs.ExchangeConfig
		errorValue error
		wantError  bool
	}{
		// Test Cases
		{},
	}
	for _, tt := range tests {
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

// Ticker & Trade Handlers
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
		r1, err := HandleTickerMessage(tt.message, tt.wrapped)

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
		r1, err := HandleTradeMessage(tt.message, tt.wrapped)

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
