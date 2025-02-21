package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Methods to add data to the buffer
func (c *DataBuffer) AddData(record interface{}) error {
	switch v := record.(type) {
	case []TickerDataStruct:
		c.tickerBuffer = append(c.tickerBuffer, v)
		if len(c.tickerBuffer) >= c.maxSize {
			return c.FlushData()
		}
	case []TradeDataStruct:
		c.tradeBuffer = append(c.tradeBuffer, v)
		if len(c.tradeBuffer) >= c.maxSize {
			return c.FlushData()
		}
	default:
		return fmt.Errorf("unsupported data type")
	}
	return nil
}

func FormatData(record interface{}) []string {
	switch v := record.(type) {
	case TickerDataStruct:
		return []string{
			strconv.FormatInt(int64(v.TimeStamp), 10),
			strconv.FormatInt(int64(v.Date), 10),
			v.Symbol,
			v.BidPrice,
			v.BidSize,
			v.AskPrice,
			v.AskSize,
		}
	case TradeDataStruct:
		return []string{
			strconv.FormatInt(int64(v.TimeStamp), 10),
			strconv.FormatInt(int64(v.Date), 10),
			v.Symbol,
			v.Price,
			v.Quantity,
			strconv.FormatBool(v.Bid_MM),
		}
	default:
		return nil
	}
}

func (c *DataBuffer) FlushData() error {
	file, err := os.OpenFile(c.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	isEmpty := fileIsEmpty(c.fileName)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header if file is empty
	if isEmpty {
		header := []string{"TimeStamp", "Date", "Symbol", "BidPrice", "BidSize", "AskPrice", "AskSize"}
		if c.dataType == "trade" {
			header = []string{"TimeStamp", "Date", "Symbol", "Price", "Quantity", "Bid_MM"}
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %w", err)
		}
	}

	// Write buffer data
	if c.dataType == "ticker" {
		for _, batch := range c.tickerBuffer {
			for _, record := range batch {
				if err := writer.Write(FormatData(record)); err != nil {
					return fmt.Errorf("error writing record to CSV: %w", err)
				}
			}
		}
		c.tickerBuffer = make([][]TickerDataStruct, 0) // Reset efficiently
	} else if c.dataType == "trade" {
		for _, batch := range c.tradeBuffer {
			for _, record := range batch {
				if err := writer.Write(FormatData(record)); err != nil {
					return fmt.Errorf("error writing record to CSV: %w", err)
				}
			}
		}
		c.tradeBuffer = make([][]TradeDataStruct, 0) // Reset efficiently
	} else {
		return fmt.Errorf("unsupported data type")
	}

	return nil
}

// Helper function to check if file is empty
func fileIsEmpty(filename string) bool {
	fileInfo, err := os.Stat(filename)
	return os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0)
}

// Create a new buffer
func NewDataBuffer(dataType string, dataStream string, maxSize int, fileName string) *DataBuffer {
	return &DataBuffer{
		tickerBuffer: make([][]TickerDataStruct, 0),
		tradeBuffer:  make([][]TradeDataStruct, 0),
		dataType:     dataType,
		dataStream:   dataStream,
		maxSize:      maxSize,
		fileName:     fileName,
	}
}
