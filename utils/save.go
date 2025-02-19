package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Method to add data to the buffer

func (c *TickerDataBuffer) AddData(record []TickerDataStruct) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.FlushData()
	}
}

func (c *TradeDataBuffer) AddData(record []TradeDataStruct) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.FlushData()
	}
}

// Method to save the buffer to a CSV file.

func (c *TickerDataBuffer) FlushData() error {
	// Check if the file exists and is empty.
	fileInfo, err := os.Stat(c.fileName)
	isEmpty := os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error checking file: %w", err)
	}

	// Open or create the CSV file for appending.
	file, err := os.OpenFile(c.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header if the file is empty.
	if isEmpty {
		header := []string{"TimeStamp", "Date", "Symbol", "BidPrice", "BidSize", "AskPrice", "AskSize"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %w", err)
		}
	}

	// Write data from buffer to CSV.
	for _, batch := range c.buffer {
		for _, record := range batch {
			stringRecord := []string{
				strconv.Itoa(int(record.TimeStamp)),
				strconv.Itoa(int(record.Date)),
				record.Symbol,
				strconv.FormatFloat(float64(record.BidPrice), 'f', 4, 32),
				strconv.FormatFloat(float64(record.BidSize), 'f', 4, 32),
				strconv.FormatFloat(float64(record.AskPrice), 'f', 4, 32),
				strconv.FormatFloat(float64(record.AskSize), 'f', 4, 32),
			}
			if err := writer.Write(stringRecord); err != nil {
				return fmt.Errorf("error writing record to CSV: %w", err)
			}
		}
	}

	// Reset buffer efficiently.
	c.buffer = nil

	return nil
}

func (c *TradeDataBuffer) FlushData() error {
	// Check if the file exists and is empty.
	fileInfo, err := os.Stat(c.fileName)
	isEmpty := os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error checking file: %w", err)
	}

	// Open or create the CSV file for appending.
	file, err := os.OpenFile(c.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header if the file is empty.
	if isEmpty {
		header := []string{"TimeStamp", "Date", "Symbol", "Price", "Quantity", "Bid_MM"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %w", err)
		}
	}

	// Write data from buffer to CSV.
	for _, batch := range c.buffer {
		for _, record := range batch {
			stringRecord := []string{
				strconv.Itoa(int(record.TimeStamp)),
				strconv.Itoa(int(record.Date)),
				record.Symbol,
				strconv.FormatFloat(float64(record.Price), 'f', 4, 32),    // Fixed precision
				strconv.FormatFloat(float64(record.Quantity), 'f', 4, 32), // Fixed precision
				strconv.FormatBool(record.Bid_MM),
			}
			if err := writer.Write(stringRecord); err != nil {
				return fmt.Errorf("error writing record to CSV: %w", err)
			}
		}
	}

	// Reset buffer efficiently.
	c.buffer = nil

	return nil
}

// Create a new buffer

func NewTickerCSVBuffer(maxSize int, fileName string) *TickerDataBuffer {
	return &TickerDataBuffer{
		buffer:   make([][]TickerDataStruct, 0),
		maxSize:  maxSize,
		fileName: fileName,
	}
}

func NewTradeCSVBuffer(maxSize int, fileName string) *TradeDataBuffer {
	return &TradeDataBuffer{
		buffer:   make([][]TradeDataStruct, 0),
		maxSize:  maxSize,
		fileName: fileName,
	}
}
