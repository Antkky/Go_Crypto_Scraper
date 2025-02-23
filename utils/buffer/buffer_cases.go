package buffer

import (
	"github.com/Antkky/go_crypto_scraper/utils"
)

var AddDataTestCases = []struct {
	name       string
	dataType   string
	data       interface{}
	errorValue string
	wantError  bool
}{
	{
		name:     "Valid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "97,242.02",
			Quantity:  "12",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  false,
	},
	{
		name:     "Valid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "97,242.02",
			Quantity:  "12",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  false,
	},
	{
		name:     "Invalid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			Date:   0,
			Symbol: "BTCUSD",
			Price:  "97,242.02",
			Bid_MM: false,
		},
		errorValue: "",
		wantError:  true,
	},
	{
		name:     "Invalid Trade Data 2",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "",
			Quantity:  "",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  true,
	},
}

var FlushDataTestCases = []struct {
	name       string
	dataType   string
	data       interface{}
	errorValue string
	wantError  bool
}{
	{
		name:     "Valid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "97,242.02",
			Quantity:  "12",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  false,
	},
	{
		name:     "Valid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "97,242.02",
			Quantity:  "12",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  false,
	},
	{
		name:     "Invalid Trade Data 1",
		dataType: "trade",
		data: utils.TradeDataStruct{
			Date:   0,
			Symbol: "BTCUSD",
			Price:  "97,242.02",
			Bid_MM: false,
		},
		errorValue: "",
		wantError:  true,
	},
	{
		name:     "Invalid Trade Data 2",
		dataType: "trade",
		data: utils.TradeDataStruct{
			TimeStamp: 1231231,
			Date:      0,
			Symbol:    "BTCUSD",
			Price:     "",
			Quantity:  "",
			Bid_MM:    false,
		},
		errorValue: "",
		wantError:  true,
	},
}
