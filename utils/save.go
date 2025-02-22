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
		c.TickerBuffer = append(c.TickerBuffer, v)
		if uint16(len(c.TickerBuffer)) >= c.MaxSize {
			return c.FlushData()
		}
	case []TradeDataStruct:
		c.TradeBuffer = append(c.TradeBuffer, v)
		if uint16(len(c.TradeBuffer)) >= c.MaxSize {
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
	filepath := fmt.Sprintf("%s%s", c.FilePath, c.FileName)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	isEmpty := fileIsEmpty(filepath)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header if file is empty
	if isEmpty {
		header := []string{"TimeStamp", "Date", "Symbol", "BidPrice", "BidSize", "AskPrice", "AskSize"}
		if c.DataType == "trade" {
			header = []string{"TimeStamp", "Date", "Symbol", "Price", "Quantity", "Bid_MM"}
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %w", err)
		}
	}

	// Write buffer data
	if c.DataType == "ticker" {
		for _, batch := range c.TickerBuffer {
			for _, record := range batch {
				if err := writer.Write(FormatData(record)); err != nil {
					return fmt.Errorf("error writing record to CSV: %w", err)
				}
			}
		}
		c.TickerBuffer = make([][]TickerDataStruct, 0)
	} else if c.DataType == "trade" {
		for _, batch := range c.TradeBuffer {
			for _, record := range batch {
				if err := writer.Write(FormatData(record)); err != nil {
					return fmt.Errorf("error writing record to CSV: %w", err)
				}
			}
		}
		c.TradeBuffer = make([][]TradeDataStruct, 0)
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
func NewDataBuffer(dataType string, market string, id string, maxSize int, fileName string, filePath string) *DataBuffer {
	return &DataBuffer{
		TickerBuffer: make([][]TickerDataStruct, 0),
		TradeBuffer:  make([][]TradeDataStruct, 0),
		DataType:     dataType,
		Market:       market,
		ID:           id,
		MaxSize:      uint16(maxSize),
		FileName:     fileName,
		FilePath:     filePath,
	}
}
