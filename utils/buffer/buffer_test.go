package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	for _, tt := range FlushDataTestCases {
		buffer := NewDataBuffer(tt.dataType, "spot", "TestAddData", 10, "Test.csv", "../../Data/Tests")
		err := buffer.AddData(tt.data)
		if err != nil {
			if tt.wantError {
				assert.EqualError(t, err, tt.errorValue)
			} else {
				assert.Error(t, err)
			}
		}

		// check the buffer for the added data
		if tt.dataType == "trade" {
			assert.Contains(t, buffer.TradeBuffer, tt.data)
		} else {
			assert.Contains(t, buffer.TickerBuffer, tt.data)
		}

		err = buffer.FlushData() // flush the data

		if err != nil {
			if tt.wantError {
				assert.EqualError(t, err, tt.errorValue)
			} else {
				assert.Error(t, err)
			}
		}
	}
}
