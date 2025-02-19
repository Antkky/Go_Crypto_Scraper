package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Buffer Structs

type TickerDataBuffer struct {
	buffer   [][]TickerDataStruct
	maxSize  int
	fileName string
}

type TradeDataBuffer struct {
	buffer   [][]TradeDataStruct
	maxSize  int
	fileName string
}

// Method to add data to the buffer
// buffer.AddData([]structs.TickerData)

func (c *TickerDataBuffer) AddData(record []TickerDataStruct) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.FlushData()
		c.buffer = nil
	}
}

func (c *TradeDataBuffer) AddData(record []TradeDataStruct) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.FlushData()
		c.buffer = nil
	}
}

// Method to save the buffer to a CSV file.
// buffer.FlushData([]structs.TickerData)

// these two functions are the same but for different structs
func (c *TickerDataBuffer) FlushData() {
	// Create or open the CSV file for appending.
	file, err := os.OpenFile(c.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the data from the buffer to the CSV.
	for _, record := range c.buffer {
		err := writer.Write(record) // this takes in a string slice
		if err != nil {
			fmt.Println("Error writing record to CSV:", err)
		}
	}
}

// these two functions are the same but for different structs
func (c *TradeDataBuffer) FlushData() {
	// Create or open the CSV file for appending.
	file, err := os.OpenFile(c.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the data from the buffer to the CSV.
	for _, record := range c.buffer {
		err := writer.Write(record) // this takes in a string slice
		if err != nil {
			fmt.Println("Error writing record to CSV:", err)
		}
	}
}

// Create a new buffer
// NewTickerCSVBuffer(100, "ticker.csv")

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
