package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Antkky/go_crypto_scraper/structs"
)

// Buffer Structs

type TickerDataBuffer struct {
	buffer   [][]structs.TickerData
	maxSize  int
	fileName string
}

type TradeDataBuffer struct {
	buffer   [][]structs.TradeData
	maxSize  int
	fileName string
}

// Method to add data to the buffer
// buffer.AddData([]structs.TickerData)

func (c *TickerDataBuffer) AddData(record []structs.TickerData) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.SaveToCSV()
		c.buffer = nil
	}
}

func (c *TradeDataBuffer) AddData(record []structs.TradeData) {
	c.buffer = append(c.buffer, record)
	if len(c.buffer) >= c.maxSize {
		c.SaveToCSV()
		c.buffer = nil
	}
}

// Method to save the buffer to a CSV file.
// buffer.AddData([]structs.TickerData)

func (c *TickerDataBuffer) SaveToCSV() {
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

func (c *TradeDataBuffer) SaveToCSV() {
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
		buffer:   make([][]structs.TickerData, 0),
		maxSize:  maxSize,
		fileName: fileName,
	}
}

func NewTradeCSVBuffer(maxSize int, fileName string) *TradeDataBuffer {
	return &TradeDataBuffer{
		buffer:   make([][]structs.TradeData, 0),
		maxSize:  maxSize,
		fileName: fileName,
	}
}
