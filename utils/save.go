package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// validateFilePath checks if the directory exists and creates it if it does not
func validateFilePath(savefilepath string) error {
	dir := filepath.Dir(savefilepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %s, error: %w", dir, err)
		}
	}
	return nil
}

// isFileEmpty checks if the file is empty or not
func isFileEmpty(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // file doesn't exist, so it's empty
		}
		return false, fmt.Errorf("error getting file info: %w", err)
	}
	return fileInfo.Size() == 0, nil
}

// getCSVHeader returns the correct CSV header based on the data type
func getCSVHeader(dataType string) ([]string, error) {
	switch dataType {
	case "trade":
		return []string{"TimeStamp", "Date", "Symbol", "Price", "Quantity", "Bid_MM"}, nil
	case "ticker":
		return []string{"TimeStamp", "Date", "Symbol", "BidPrice", "BidSize", "AskPrice", "AskSize"}, nil
	default:
		return nil, fmt.Errorf("unsupported data type for header: %s", dataType)
	}
}

// writeDataToCSV writes a batch of data to the CSV file
func writeDataToCSV(writer *csv.Writer, buffer interface{}) error {
	switch batch := buffer.(type) {
	case []TickerDataStruct:
		for _, record := range batch {
			fData, err := FormatData(record)
			if err != nil {
				return errors.New("error formatting data")
			}
			if err = writer.Write(fData); err != nil {
				return fmt.Errorf("error writing ticker record: %w", err)
			}
		}
	case []TradeDataStruct:
		for _, record := range batch {
			fData, err := FormatData(record)
			if err != nil {
				return fmt.Errorf("error formatting data: %s", err)
			}
			if err = writer.Write(fData); err != nil {
				return fmt.Errorf("error writing trade record: %w", err)
			}
		}
	default:
		return fmt.Errorf("unsupported buffer type")
	}
	return nil
}

// FormatData formats the data for writing to CSV
func FormatData(record interface{}) ([]string, error) {
	switch v := record.(type) {
	case TickerDataStruct:
		if v.BidPrice == "" || v.AskPrice == "" || v.BidSize == "" || v.AskSize == "" {
			return nil, fmt.Errorf("missing required field(s) in TickerDataStruct")
		}

		return []string{
			fmt.Sprintf("%d", v.TimeStamp),
			fmt.Sprintf("%d", v.Date),
			v.Symbol,
			v.BidPrice,
			v.BidSize,
			v.AskPrice,
			v.AskSize,
		}, nil

	case TradeDataStruct:
		// Check for empty fields that should contain data
		if v.Price == "" || v.Quantity == "" {
			return nil, nil
		}

		// Convert timestamp and date to string
		return []string{
			fmt.Sprintf("%d", v.TimeStamp), // TimeStamp as integer
			fmt.Sprintf("%d", v.Date),      // Date as integer
			v.Symbol,                       // Symbol
			v.Price,                        // Price as float with 2 decimals
			v.Quantity,                     // Quantity as integer
			fmt.Sprintf("%t", v.Bid_MM),    // Bid_MM as string ("true" or "false")
		}, nil

	default:
		return nil, fmt.Errorf("unsupported record type: %T", record)
	}
}

// Methods to add data to the buffer
func (c *DataBuffer) AddData(record interface{}) error {
	switch v := record.(type) {
	case TickerDataStruct:
		c.TickerBuffer = append(c.TickerBuffer, v)
		if len(c.TickerBuffer) >= c.MaxSize {
			return c.FlushData()
		}
	case TradeDataStruct:
		c.TradeBuffer = append(c.TradeBuffer, v)
		if len(c.TradeBuffer) >= c.MaxSize {
			return c.FlushData()
		}
	default:
		return fmt.Errorf("unsupported data type")
	}
	return nil
}

func (c *DataBuffer) FlushData() error {
	filepath := fmt.Sprintf("%s/%s", c.FilePath, c.FileName)
	if err := validateFilePath(filepath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	log.Println("Writing to file:", filepath)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	isEmpty, err := isFileEmpty(filepath)
	if err != nil {
		return fmt.Errorf("error checking file empty status: %w", err)
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if isEmpty {
		header, err := getCSVHeader(c.DataType)
		if err != nil {
			return fmt.Errorf("error getting CSV header: %w", err)
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %w", err)
		}
	}
	if c.DataType == "ticker" {
		if err := writeDataToCSV(writer, c.TickerBuffer); err != nil {
			return fmt.Errorf("error writing ticker records to CSV: %w", err)
		}
		c.TickerBuffer = nil
	} else if c.DataType == "trade" {
		if err := writeDataToCSV(writer, c.TradeBuffer); err != nil {
			return fmt.Errorf("error writing trade records to CSV: %w", err)
		}
		c.TradeBuffer = nil
	} else {
		return fmt.Errorf("unsupported data type: %s", c.DataType)
	}
	return nil
}

// Create a new buffer
func NewDataBuffer(dataType string, market string, id string, maxSize int, fileName string, filePath string) *DataBuffer {
	return &DataBuffer{
		TickerBuffer: make([]TickerDataStruct, maxSize),
		TradeBuffer:  make([]TradeDataStruct, maxSize),
		DataType:     dataType,
		Market:       market,
		ID:           id,
		MaxSize:      maxSize,
		FileName:     fileName,
		FilePath:     filePath,
	}
}
