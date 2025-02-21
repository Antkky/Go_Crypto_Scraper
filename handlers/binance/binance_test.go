package binance

import (
	"testing"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/stretchr/testify/assert"
)

// TestProcessMessage
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
func TestProcessMessage(t *testing.T) {
	for _, tt := range ProcessMessageTypeCases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				r1 utils.TickerDataStruct
				r2 utils.TradeDataStruct
			)

			// Call the ProcessMessage function
			err := ProcessMessage(tt.message, &r1, &r2)

			// Error handling logic
			if tt.wantError {
				// Assert that the error is not nil and matches the expected error
				assert.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tt.errorValue, "Error type does not match expected")
			} else {
				// Ensure no error occurred if not expected
				assert.NoError(t, err, "Unexpected error occurred")
			}

			// Validate the result based on event type
			switch tt.eventType {
			case "24hrTicker":
				assert.Equal(t, tt.r1, r1, "Ticker data (r1) does not match expected output")
			case "trade":
				assert.Equal(t, tt.r2, r2, "Trade data (r2) does not match expected output")
			default:
				t.Errorf("Unexpected event type: %s", tt.eventType)
			}
		})
	}
}

// TestHandleMessage
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
			if tt.wantError {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, tt.errorValue, "Error type does not match expected")
				}
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}
